package monitor

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/jbrukh/window"
	"github.com/samalba/dockerclient"

	"github.com/elleFlorio/gru/service"
)

var (
	monitorActive  bool
	gruStats       GruStats
	history        statsHistory
	ErrNoIndexById error = errors.New("No index for such Id")
)

const W_SIZE = 50
const W_MULT = 1000

type monitor struct {
	c_stop chan struct{}
	c_err  chan error
}

func init() {
	gruStats = GruStats{
		Service:  make(map[string]ServiceStats),
		Instance: make(map[string]InstanceStats),
	}

	history = statsHistory{make(map[string]instanceHistory)}
}

func NewMonitor(c_stop chan struct{}, c_err chan error) *monitor {
	return &monitor{
		c_stop,
		c_err,
	}
}

func (p *monitor) Run() GruStats {
	snapshot := GruStats{
		Service:  make(map[string]ServiceStats),
		Instance: make(map[string]InstanceStats),
	}

	makeSnapshot(&gruStats, &snapshot)

	services := service.List()
	for _, name := range services {
		resetEventsStats(name, &gruStats)
	}

	return snapshot
}

func makeSnapshot(src *GruStats, dst *GruStats) {
	// Copy service stats
	for k, v := range src.Service {
		srv_src := v
		// Copy instances status
		status_src := v.Instances
		all_dst := make([]string, len(status_src.All), len(status_src.All))
		runnig_dst := make([]string, len(status_src.Running), len(status_src.Running))
		stopped_dst := make([]string, len(status_src.Stopped), len(status_src.Stopped))
		paused_dst := make([]string, len(status_src.Paused), len(status_src.Paused))
		copy(all_dst, status_src.All)
		copy(runnig_dst, status_src.Running)
		copy(stopped_dst, status_src.Stopped)
		copy(paused_dst, status_src.Paused)
		status_dst := InstanceStatus{
			all_dst,
			runnig_dst,
			stopped_dst,
			paused_dst,
		}
		// Copy events
		events_src := srv_src.Events
		stop_dst := make([]string, len(events_src.Stop), len(events_src.Stop))
		start_dst := make([]string, len(events_src.Start), len(events_src.Start))
		copy(start_dst, events_src.Start)
		copy(stop_dst, events_src.Stop)
		// Create new service stats
		events_dst := EventStats{
			start_dst,
			stop_dst,
		}
		srv_dst := ServiceStats{status_dst, events_dst}
		dst.Service[k] = srv_dst
	}

	//Copy instance stats

	for k, v := range history.instance {
		instCpuHist := v.cpu.totalUsage.Slice()
		instCpuSysHist := v.cpu.sysUsage.Slice()
		instCpu_dst := make([]float64, len(instCpuHist), len(instCpuHist))
		copy(instCpu_dst, instCpuHist)
		instCpuSys_dst := make([]float64, len(instCpuSysHist), len(instCpuSysHist))
		copy(instCpuSys_dst, instCpuSysHist)
		cpuStats_dst := CpuStats{
			TotalUsage: instCpu_dst,
			SysUsage:   instCpuSys_dst,
		}
		instStats_dst := InstanceStats{
			Cpu: cpuStats_dst,
		}
		dst.Instance[k] = instStats_dst
	}

	//Copy system stats
	//none for now...
}

func resetEventsStats(srvName string, stats *GruStats) {
	srvStats := stats.Service[srvName]

	log.WithFields(log.Fields{
		"status":  "monitored events",
		"service": srvName,
		"start":   srvStats.Events.Start,
		"stop":    srvStats.Events.Stop,
	}).Debugln("Running monitor")

	srvStats.Events = EventStats{}
	stats.Service[srvName] = srvStats
}

func (p *monitor) Start(docker *dockerclient.DockerClient) {
	log.WithField("status", "start").Debugln("Running monitor")
	monitorActive = true
	c_mntrerr := make(chan error)
	c_evntstart := make(chan string)

	docker.StartMonitorEvents(eventCallback, c_mntrerr, c_evntstart)

	// Get the list of containers (running or not) to monitor
	containers, err := docker.ListContainers(true, false, "")
	if err != nil {
		p.monitorError(err)
	}

	// Start the monitor for each configured service
	for _, c := range containers {
		info, _ := docker.InspectContainer(c.Id)
		status := getContainerStatus(info)
		srv, err := service.GetServiceByImage(c.Image)
		if err != nil {
			log.WithFields(log.Fields{
				"err":   err,
				"image": c.Image,
			}).Warningln("Running monitor")
		} else {
			addResource(c.Id, srv.Name, status, &gruStats, &history)
			docker.StartMonitorStats(c.Id, statCallBack, c_mntrerr)
		}
	}

	for monitorActive {
		select {
		case err := <-c_mntrerr:
			p.monitorError(err)
		case newId := <-c_evntstart:
			docker.StartMonitorStats(newId, statCallBack, c_mntrerr)
		}
	}
}

// Events are: create, destroy, die, exec_create, exec_start, export, kill, oom, pause, restart, start, stop, unpause
func eventCallback(event *dockerclient.Event, ec chan error, args ...interface{}) {
	log.WithFields(log.Fields{
		"status": "received event",
		"event":  event.Status,
		"from":   event.From,
	}).Debug("Running monitor")

	c_evntstart := args[0].(chan string)

	switch event.Status {
	case "stop":
	case "die":
		removeResource(event.Id, &gruStats, &history)
	case "create":
	case "start":
		// TODO handle error
		srv, err := service.GetServiceByImage(event.From)
		if err != nil {
			log.WithFields(log.Fields{
				"status": "resource not added",
				"error":  err,
			}).Warnln("Running monitor")
		} else {
			addResource(event.Id, srv.Name, "running", &gruStats, &history)
			c_evntstart <- event.Id
		}

	default:
		log.WithFields(log.Fields{
			"status": "event not handled",
			"event":  event.Status,
			"from":   event.From,
		}).Warnln("Running monitor")
	}

}

func getContainerStatus(info *dockerclient.ContainerInfo) string {
	if info.State.Running {
		return "running"
	} else if info.State.Paused {
		return "paused"
	} else {
		return "stopped"
	}
}

func addResource(id string, srvName string, status string, stats *GruStats, hist *statsHistory) {
	servStats := stats.Service[srvName]
	_, err := findIdIndex(id, servStats.Instances.All)
	if err != nil {
		servStats.Instances.All = append(servStats.Instances.All, id)
	}

	switch status {
	case "running":
		servStats.Instances.Running = append(servStats.Instances.Running, id)
		index, err := findIdIndex(id, servStats.Instances.Stopped)
		if err == nil {
			servStats.Instances.Stopped = append(
				servStats.Instances.Stopped[:index],
				servStats.Instances.Stopped[index+1:]...)
		}

		servStats.Events.Start = append(servStats.Events.Start, id)
		cpu := cpuHistory{
			totalUsage: window.New(W_SIZE, W_MULT),
			sysUsage:   window.New(W_SIZE, W_MULT),
		}

		hist.instance[id] = instanceHistory{
			cpu: cpu,
		}
	case "stopped":
		servStats.Instances.Stopped = append(servStats.Instances.Stopped, id)
	case "paused":
		servStats.Instances.Paused = append(servStats.Instances.Paused, id)
	default:
		log.WithFields(log.Fields{
			"error":   "Unknown container state: " + status,
			"service": srvName,
			"id":      id,
		}).Warnln("Running monitor")
	}
	stats.Service[srvName] = servStats

	log.WithFields(log.Fields{
		"status":  "started monitor on container",
		"service": srvName,
		"id":      id,
	}).Infoln("Running monitor")
}

func removeResource(id string, stats *GruStats, hist *statsHistory) {
	srvName := findServiceByInstanceId(id, stats)

	// Updating service stats
	srvStats := stats.Service[srvName]
	running := srvStats.Instances.Running
	index, _ := findIdIndex(id, running)
	running = append(running[:index], running[index+1:]...)
	srvStats.Instances.Running = running
	srvStats.Instances.Stopped = append(srvStats.Instances.Stopped, id)

	// Upating Event stats
	srvStats.Events.Stop = append(srvStats.Events.Stop, id)

	stats.Service[srvName] = srvStats

	// Updating Instances stats
	delete(stats.Instance, id)

	delete(hist.instance, id)

	log.WithFields(log.Fields{
		"status":  "removed instance",
		"service": srvName,
		"id":      id,
	}).Infoln("Running monitor")
}

// TODO create error?
func findServiceByInstanceId(id string, stats *GruStats) string {
	for k, v := range stats.Service {
		for _, instance := range v.Instances.All {
			if instance == id {
				return k
			}
		}
	}

	return ""
}

func findIdIndex(id string, instances []string) (int, error) {
	for index, v := range instances {
		if v == id {
			return index, nil
		}
	}

	return -1, ErrNoIndexById
}

func statCallBack(id string, stats *dockerclient.Stats, ec chan error, args ...interface{}) {
	instHist := history.instance[id]

	// Instance stats update

	// Cpu history usage update
	totCpu := float64(stats.CpuStats.CpuUsage.TotalUsage)
	sysCpu := float64(stats.CpuStats.SystemUsage)
	instHist.cpu.totalUsage.PushBack(totCpu)
	instHist.cpu.sysUsage.PushBack(sysCpu)
	history.instance[id] = instHist
}

func (p *monitor) monitorError(err error) {
	log.WithFields(log.Fields{
		"status": "error",
		"error:": err,
	}).Errorln("Running monitor")
	p.c_err <- err
}

func (p *monitor) Stop() {
	monitorActive = false
	log.WithField("status", "done").Warnln("Running monitor")
	p.c_stop <- struct{}{}
}
