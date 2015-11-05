package monitor

import (
	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/jbrukh/window"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/samalba/dockerclient"

	"github.com/elleFlorio/gru/autonomic/monitor/metric"
	"github.com/elleFlorio/gru/container"
	"github.com/elleFlorio/gru/service"
)

var optionsLog = dockerclient.LogOptions{Follow: true, Stdout: true, Stderr: true, Tail: 1}

func Start(cError chan error, cStop chan struct{}) {
	log.WithField("status", "start").Debugln("Autonomic Monitor")
	metric.Manager().Start()

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

		// ###########################################################
		// NOT USED BY NOW. I keep it for now because I'm not sure that
		// what I'm doing makes sense...

		// This is related to the service, not the single isntance,
		// so I have to check if it's already initialized.
		// Maybe I can find a better way to do this...
		/*if _, ok := hist.service[srvName]; !ok {
			respTime := window.New(W_SIZE, W_MULT)
			hist.service[srvName] = metricsHistory{respTime}
		}*/
		// ###########################################################

		hist.instance[id] = instanceHistory{cpu}

		contLog, err := container.Docker().Client.ContainerLogs(id, &optionsLog)
		if err != nil {
			log.WithField("error", err).Errorln("Cannot start log monitoring on container ", id)
		} else {
			metric.Manager().StartCollector(contLog)
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
		}).Warnln("Cannot add resource to monitor")
	}
	stats.Service[srvName] = servStats

	log.WithFields(log.Fields{
		"status":  status,
		"service": srvName,
	}).Infoln("Added resource to monitor")
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
		index, err = findIdIndex(id, pending)
		if err != nil {
			log.WithField("id", id).Errorln("Cannot find pending container to stop")
			return
		}
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
	// FIXME this can be a problem if the instance is killed
	// While I'm computing the stats to make the snapshot
	// FIX this issue
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
