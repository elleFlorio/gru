package monitor

import (
	log "github.com/Sirupsen/logrus"
	"github.com/samalba/dockerclient"

	"github.com/elleFlorio/gru/service"
)

var monitorActive bool
var gruStats GruStats

type monitor struct {
	c_stop chan struct{}
	c_err  chan error
}

func init() {
	gruStats = GruStats{
		Service:  make(map[string]ServiceStats),
		Instance: make(map[string]InstanceStats),
	}
}

func NewMonitor(c_stop chan struct{}, c_err chan error) *monitor {
	return &monitor{
		c_stop,
		c_err,
	}
}

// FIXME need to find a way to reset events
func (p *monitor) Run() GruStats {
	updGruStats := GruStats{
		Service:  make(map[string]ServiceStats),
		Instance: make(map[string]InstanceStats),
	}

	copyStats(&gruStats, &updGruStats)

	services := service.List()
	for _, name := range services {
		resetEventsStats(name, &gruStats)
	}

	return updGruStats
}

func copyStats(src *GruStats, dst *GruStats) {
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
		copy(start_dst, events_src.Stop)
		copy(stop_dst, events_src.Stop)
		// Create new service stats
		events_dst := EventStats{
			start_dst,
			stop_dst,
		}
		srv_dst := ServiceStats{status_dst, events_dst}
		dst.Service[k] = srv_dst
	}

	for k, v := range src.Instance {
		dst.Instance[k] = v
	}

	dst.System.Cpu = src.System.Cpu
}

func resetEventsStats(srvName string, stats *GruStats) {
	srvStats := stats.Service[srvName]
	srvStats.Events = EventStats{}
	stats.Service[srvName] = srvStats
}

func (p *monitor) Start(docker *dockerclient.DockerClient) {
	log.WithField("status", "start").Debugln("Running monitor")
	monitorActive = true
	c_mntrerr := make(chan error)

	docker.StartMonitorEvents(eventCallback, c_mntrerr)

	// Get the list of active containers to monitor
	containers, err := docker.ListContainers(false, false, "")
	if err != nil {
		p.monitorError(err)
	}

	// Start the monitor for each active container
	for _, c := range containers {
		serv, err := service.GetServiceByImage(c.Image)
		if err != nil {
			p.monitorError(err)
		} else {
			servStats := gruStats.Service[serv.Name]
			servStats.Instances.All = append(servStats.Instances.All, c.Id)
			servStats.Events.Start = append(servStats.Events.Start, c.Id)
			servStats.Instances.Running = append(servStats.Instances.Running, c.Id)
			gruStats.Service[serv.Name] = servStats
			docker.StartMonitorStats(c.Id, statCallBack, c_mntrerr)

			log.WithFields(log.Fields{
				"id":    c.Id,
				"image": c.Image,
			}).Infoln("Started monitor on container")
		}
	}

	for monitorActive {
		select {
		case err := <-c_mntrerr:
			p.monitorError(err)
		}
	}
}

// Events are: create, destroy, die, exec_create, exec_start, export, kill, oom, pause, restart, start, stop, unpause
func eventCallback(event *dockerclient.Event, ec chan error, args ...interface{}) {
	log.WithFields(log.Fields{
		"from":   event.From,
		"status": event.Status,
	}).Debug("Received event")

	switch event.Status {
	case "stop":
	case "die":
		removeResource(event.Id, &gruStats)

		log.WithFields(log.Fields{
			"status": "removed instance",
			"image":  event.From,
			"id":     event.Id,
		}).Debug("Running monitor")

	default:
		log.WithFields(log.Fields{
			"status": "event not handled",
			"event":  event.Status,
			"from":   event.From,
		}).Warnln("Running monitor")
	}

}

func removeResource(id string, stats *GruStats) {
	srvName := findServiceByInstanceId(id, stats)

	// Updating service stats
	srvStats := stats.Service[srvName]
	running := srvStats.Instances.Running
	index := findIdIndex(id, running)
	running = append(running[:index], running[index+1:]...)
	srvStats.Instances.Running = running
	srvStats.Instances.Stopped = append(srvStats.Instances.Stopped, id)

	// Upating Event stats
	srvStats.Events.Stop = append(srvStats.Events.Stop, id)

	stats.Service[srvName] = srvStats

	// Updating Instances stats
	delete(stats.Instance, id)
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

// TODO create error?
func findIdIndex(id string, instances []string) int {
	for index, v := range instances {
		if v == id {
			return index
		}
	}

	return -1
}

func statCallBack(id string, stats *dockerclient.Stats, ec chan error, args ...interface{}) {
	InstanceStats := gruStats.Instance[id]
	InstanceStats.Cpu = stats.CpuStats.CpuUsage.TotalUsage
	gruStats.Instance[id] = InstanceStats
	gruStats.System.Cpu = stats.CpuStats.SystemUsage

	// log.WithFields(log.Fields{
	// 	"status": "update",
	// 	"id:":    id,
	// }).Debugln("Running monitor")
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
