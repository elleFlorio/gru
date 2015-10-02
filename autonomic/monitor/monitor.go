package monitor

import (
	"errors"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/jbrukh/window"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/samalba/dockerclient"

	"github.com/elleFlorio/gru/container"
	"github.com/elleFlorio/gru/node"
	"github.com/elleFlorio/gru/service"
	"github.com/elleFlorio/gru/storage"
)

var (
	monitorActive  bool
	gruStats       GruStats
	history        statsHistory
	c_stop         chan struct{}
	c_err          chan error
	ErrNoIndexById error = errors.New("No index for such Id")
)

//History window
const W_SIZE = 50
const W_MULT = 1000

//My data type
const dataType string = "stats"

func init() {
	gruStats = GruStats{
		Service:  make(map[string]ServiceStats),
		Instance: make(map[string]InstanceStats),
	}

	history = statsHistory{make(map[string]instanceHistory)}
}

func Run() {
	snapshot := GruStats{
		Service:  make(map[string]ServiceStats),
		Instance: make(map[string]InstanceStats),
	}

	for name, _ := range gruStats.Service {
		updateRunningInstances(name, &gruStats, W_SIZE)
		computeServiceCpuPerc(name, &gruStats)
	}
	computeSystemCpu(&gruStats)

	makeSnapshot(&gruStats, &snapshot)
	data, err := convertStatsToData(snapshot)
	if err != nil {
		log.WithField("error", "Cannot convert stats to data").Debugln("Running monitor")
	} else {
		storage.DataStore().StoreData(node.Config().UUID, data, "stats")
	}

	services := service.List()
	for _, name := range services {
		resetEventsStats(name, &gruStats)
	}
}

func updateRunningInstances(name string, stats *GruStats, wsize int) {
	srvStats := stats.Service[name]
	pending := srvStats.Instances.Pending
	for _, inst := range pending {
		if history.instance[inst].cpu.sysUsage.Size() >= W_SIZE {
			addResource(inst, name, "running", stats, &history)
		}
	}
}

func computeServiceCpuPerc(name string, stats *GruStats) {
	sum := 0.0
	avg := 0.0
	srvStats := stats.Service[name]

	for _, id := range srvStats.Instances.Running {
		instCpus := history.instance[id].cpu.totalUsage.Slice()
		sysCpus := history.instance[id].cpu.sysUsage.Slice()
		instCpuAvg := computeInstanceCpuPerc(instCpus, sysCpus)
		inst := stats.Instance[id]
		inst.Cpu = instCpuAvg
		stats.Instance[id] = inst
		sum += instCpuAvg
	}

	avg = sum / float64(len(srvStats.Instances.Running))

	srvStats.Cpu.Avg = avg
	srvStats.Cpu.Tot = sum
	stats.Service[name] = srvStats
}

// Since linux compute the cpu usage in units of jiffies, it needs to be converted
// in % using the formula used in this function.
// Explaination: http://stackoverflow.com/questions/1420426/calculating-cpu-usage-of-a-process-in-linux
func computeInstanceCpuPerc(instCpus []float64, sysCpus []float64) float64 {
	sum := 0.0
	instNext := 0.0
	sysNext := 0.0
	instPrev := 0.0
	sysPrev := 0.0
	cpu := 0.0

	for i := 1; i < len(instCpus); i++ {
		instPrev = instCpus[i-1]
		sysPrev = sysCpus[i-1]
		instNext = instCpus[i]
		sysNext = sysCpus[i]
		instDelta := instNext - instPrev
		sysDelta := sysNext - sysPrev
		if sysDelta == 0 {
			cpu = 0
		} else {
			// "100 * cpu" should produce values in [0, 100]
			cpu = instDelta / sysDelta
		}
		sum += cpu
	}
	return sum / float64(len(instCpus)-1)
}

func computeSystemCpu(stats *GruStats) {
	sum := 0.0
	for _, value := range stats.Service {
		sum += value.Cpu.Tot
	}
	stats.System.Cpu = sum
}

//TODO maybe I can just compute historical data without make a deep copy
// since now I'm serializing the structure in a string of bytes...
// check this possibility.
func makeSnapshot(src *GruStats, dst *GruStats) {
	// Copy service stats
	for name, stats := range src.Service {
		srv_src := stats
		// Copy instances status
		status_src := stats.Instances
		all_dst := make([]string, len(status_src.All), len(status_src.All))
		runnig_dst := make([]string, len(status_src.Running), len(status_src.Running))
		pending_dst := make([]string, len(status_src.Pending), len(status_src.Pending))
		stopped_dst := make([]string, len(status_src.Stopped), len(status_src.Stopped))
		paused_dst := make([]string, len(status_src.Paused), len(status_src.Paused))
		copy(all_dst, status_src.All)
		copy(runnig_dst, status_src.Running)
		copy(pending_dst, status_src.Pending)
		copy(stopped_dst, status_src.Stopped)
		copy(paused_dst, status_src.Paused)
		status_dst := InstanceStatus{
			all_dst,
			runnig_dst,
			pending_dst,
			stopped_dst,
			paused_dst,
		}
		// Copy events (NEEDED?)
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
		cpu_dst := CpuStats{stats.Cpu.Avg, stats.Cpu.Tot}
		srv_dst := ServiceStats{status_dst, events_dst, cpu_dst}
		dst.Service[name] = srv_dst
	}

	//Copy instance stats
	for id, value := range src.Instance {
		inst_dst := InstanceStats{value.Cpu}
		dst.Instance[id] = inst_dst
	}

	// PROBABLY NOT NEEDED ANYMORE
	// for k, v := range history.instance {
	// 	instCpuHist := v.cpu.totalUsage.Slice()
	// 	instCpuSysHist := v.cpu.sysUsage.Slice()
	// 	instCpu_dst := make([]float64, len(instCpuHist), len(instCpuHist))
	// 	copy(instCpu_dst, instCpuHist)
	// 	instCpuSys_dst := make([]float64, len(instCpuSysHist), len(instCpuSysHist))
	// 	copy(instCpuSys_dst, instCpuSysHist)
	// 	cpuStats_dst := CpuStats{
	// 		TotalUsage: instCpu_dst,
	// 		SysUsage:   instCpuSys_dst,
	// 	}
	// 	instStats_dst := InstanceStats{
	// 		Cpu: cpuStats_dst,
	// 	}
	// 	dst.Instance[k] = instStats_dst
	// }

	//Copy system stats
	sys_status_src := src.System.Instances
	sys_all_dst := make([]string, len(sys_status_src.All), len(sys_status_src.All))
	sys_runnig_dst := make([]string, len(sys_status_src.Running), len(sys_status_src.Running))
	sys_pending_dst := make([]string, len(sys_status_src.Pending), len(sys_status_src.Pending))
	sys_stopped_dst := make([]string, len(sys_status_src.Stopped), len(sys_status_src.Stopped))
	sys_paused_dst := make([]string, len(sys_status_src.Paused), len(sys_status_src.Paused))
	copy(sys_all_dst, sys_status_src.All)
	copy(sys_runnig_dst, sys_status_src.Running)
	copy(sys_pending_dst, sys_status_src.Pending)
	copy(sys_stopped_dst, sys_status_src.Stopped)
	copy(sys_paused_dst, sys_status_src.Paused)
	sys_status_dst := InstanceStatus{
		sys_all_dst,
		sys_runnig_dst,
		sys_pending_dst,
		sys_stopped_dst,
		sys_paused_dst,
	}
	dst.System.Instances = sys_status_dst
	dst.System.Cpu = src.System.Cpu
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

func Start(cError chan error, cStop chan struct{}) {
	log.WithField("status", "start").Debugln("Running monitor")
	monitorActive = true
	c_err = cError
	cStop = cStop

	c_mntrerr := make(chan error)
	c_evntstart := make(chan string)

	container.Docker().Client.StartMonitorEvents(eventCallback, c_mntrerr, c_evntstart)

	// Get the list of containers (running or not) to monitor
	containers, err := container.Docker().Client.ListContainers(true, false, "")
	if err != nil {
		monitorError(err)
	}

	// Start the monitor for each configured service
	for _, c := range containers {
		info, _ := container.Docker().Client.InspectContainer(c.Id)
		status := getContainerStatus(info)
		srv, err := service.GetServiceByImage(c.Image)
		if err != nil {
			log.WithFields(log.Fields{
				"err":   err,
				"image": c.Image,
			}).Warningln("Running monitor")
		} else {
			// This is needed becasuse on start I don't have data
			// to analyze the container, so every running container
			// is in a pending state
			if status == "running" {
				status = "pending"
			}
			addResource(c.Id, srv.Name, status, &gruStats, &history)
			container.Docker().Client.StartMonitorStats(c.Id, statCallBack, c_mntrerr)
		}
	}

	for monitorActive {
		select {
		case err := <-c_mntrerr:
			monitorError(err)
		case newId := <-c_evntstart:
			container.Docker().Client.StartMonitorStats(newId, statCallBack, c_mntrerr)
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
			addResource(event.Id, srv.Name, "pending", &gruStats, &history)
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
		stats.System.Instances.All = append(stats.System.Instances.All, id)
	}

	switch status {
	case "running":
		index, err := findIdIndex(id, servStats.Instances.Pending)
		servStats.Instances.Running = append(servStats.Instances.Running, id)
		stats.System.Instances.Running = append(stats.System.Instances.Running, id)
		if err != nil {
			log.WithField("error", err).Errorln("Cannot find pending instance to promote running")
		} else {
			servStats.Instances.Pending = append(
				servStats.Instances.Pending[:index],
				servStats.Instances.Pending[index+1:]...)

			sysIndex, _ := findIdIndex(id, stats.System.Instances.Pending)
			stats.System.Instances.Pending = append(
				stats.System.Instances.Pending[:sysIndex],
				stats.System.Instances.Pending[sysIndex+1:]...)
		}
	case "pending":
		servStats.Instances.Pending = append(servStats.Instances.Pending, id)
		stats.System.Instances.Pending = append(stats.System.Instances.Pending, id)

		index, err := findIdIndex(id, servStats.Instances.Stopped)
		if err != nil {
			log.WithField("error", err).Warnln("Cannot find stopped instance to promote pending")
		} else {
			servStats.Instances.Stopped = append(
				servStats.Instances.Stopped[:index],
				servStats.Instances.Stopped[index+1:]...)

			sysIndex, _ := findIdIndex(id, stats.System.Instances.Stopped)
			stats.System.Instances.Stopped = append(
				stats.System.Instances.Stopped[:sysIndex],
				stats.System.Instances.Stopped[sysIndex+1:]...)
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
		stats.System.Instances.Stopped = append(stats.System.Instances.Stopped, id)
	case "paused":
		servStats.Instances.Paused = append(servStats.Instances.Paused, id)
		stats.System.Instances.Paused = append(stats.System.Instances.Paused, id)
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
	pending := srvStats.Instances.Pending

	index, err := findIdIndex(id, running)
	if err != nil {
		// If it is not runnig it should be pending
		index, _ = findIdIndex(id, pending)
		pending = append(pending[:index], pending[index+1:]...)
		srvStats.Instances.Pending = pending

		// Updating system stats
		sysIndex, _ := findIdIndex(id, stats.System.Instances.Pending)
		stats.System.Instances.Pending = append(
			stats.System.Instances.Pending[:sysIndex],
			stats.System.Instances.Pending[sysIndex+1:]...)
	} else {
		running = append(running[:index], running[index+1:]...)
		srvStats.Instances.Running = running

		// Updating system stats
		sysIndex, _ := findIdIndex(id, stats.System.Instances.Running)
		stats.System.Instances.Running = append(
			stats.System.Instances.Running[:sysIndex],
			stats.System.Instances.Running[sysIndex+1:]...)
	}

	srvStats.Instances.Stopped = append(srvStats.Instances.Stopped, id)
	stats.System.Instances.Stopped = append(stats.System.Instances.Stopped, id)

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

func monitorError(err error) {
	log.WithFields(log.Fields{
		"status": "error",
		"error:": err,
	}).Errorln("Running monitor")
	c_err <- err
}

func Stop() {
	monitorActive = false
	log.WithField("status", "done").Warnln("Running monitor")
	c_stop <- struct{}{}
}
