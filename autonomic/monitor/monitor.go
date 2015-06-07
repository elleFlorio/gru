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

func (p *monitor) Run() GruStats {
	return gruStats
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
			servStats.Instances = append(servStats.Instances, c.Id)
			gruStats.Service[serv.Name] = servStats
			docker.StartMonitorStats(c.Id, statCallBack, c_mntrerr)

			log.WithFields(log.Fields{
				"id":    c.Id,
				"image": c.Image,
			}).Debug("Started monitor on container")
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
	log.WithField("status", "done").Debugln("Running monitor")
	p.c_stop <- struct{}{}
}
