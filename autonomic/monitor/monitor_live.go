package monitor

import (
	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/jbrukh/window"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/samalba/dockerclient"

	"github.com/elleFlorio/gru/autonomic/monitor/logreader"
	"github.com/elleFlorio/gru/container"
	res "github.com/elleFlorio/gru/resources"
	"github.com/elleFlorio/gru/service"
)

//History window
const W_SIZE = 100
const W_MULT = 1000

func Start(cError chan error, cStop chan struct{}) {
	go startMonitoring(cError, cStop)
}

func startMonitoring(cError chan error, cStop chan struct{}) {
	log.Debugln("Running autonomic monitoring")
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
			}).Warningln("Error monitoring service")
		} else {
			// This is needed becasuse on start I don't have data
			// to analyze the container, so every running container
			// is in a pending state
			if status == "running" {
				status = "pending"
			}
			addResource(c.Id, srv.Name, status, &gruStats, &history)
			container.Docker().Client.StartMonitorStats(c.Id, statCallBack, c_mntrerr)
			if status == "pending" {
				startMonitorLog(c.Id)
			}
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
	c_evntstart := args[0].(chan string)

	switch event.Status {
	case "create":
		log.WithField("image", event.From).Debugln("Created new container")
	case "start":
		// TODO handle error
		srv, err := service.GetServiceByImage(event.From)
		if err != nil {
			log.WithField("err", err).Warnln("Cannot add resource")
		} else {
			addResource(event.Id, srv.Name, "pending", &gruStats, &history)
			startMonitorLog(event.Id)
			c_evntstart <- event.Id
		}
	case "stop":
		log.WithField("image", event.From).Debugln("Received stop signal")
	case "kill":
		log.WithField("image", event.From).Debugln("Received kill signal")
	case "die":
		removeResource(event.Id, &gruStats, &history)
	default:
		log.WithFields(log.Fields{
			"err":   "event not handled",
			"event": event.Status,
			"image": event.From,
		}).Debugln("Received unknown signal")
	}

}

func addResource(id string, srvName string, status string, stats *GruStats, hist *statsHistory) {
	srv, _ := service.GetServiceByName(srvName)

	_, err := findIdIndex(id, srv.Instances.All)
	if err != nil {
		srv.Instances.All = append(srv.Instances.All, id)
	}

	switch status {
	case "running":
		index, err := findIdIndex(id, srv.Instances.Pending)
		srv.Instances.Running = append(srv.Instances.Running, id)
		if err != nil {
			log.WithField("error", err).Errorln("Cannot find pending instance to promote running")
		} else {
			srv.Instances.Pending = append(
				srv.Instances.Pending[:index],
				srv.Instances.Pending[index+1:]...)
		}
	case "pending":
		srv.Instances.Pending = append(srv.Instances.Pending, id)
		index, err := findIdIndex(id, srv.Instances.Stopped)
		if err != nil {
			log.WithField("error", err).Debugln("Cannot find stopped instance to promote pending")
		} else {
			srv.Instances.Stopped = append(
				srv.Instances.Stopped[:index],
				srv.Instances.Stopped[index+1:]...)
		}

		servStats := stats.Service[srvName]
		servStats.Events.Start = append(servStats.Events.Start, id)
		stats.Service[srvName] = servStats

		cpu := cpuHistory{
			totalUsage: window.New(W_SIZE, W_MULT),
			sysUsage:   window.New(W_SIZE, W_MULT),
		}
		mem := window.New(W_SIZE, W_MULT)
		hist.instance[id] = instanceHistory{cpu, mem}
	case "stopped":
		srv.Instances.Stopped = append(srv.Instances.Stopped, id)
		log.Debugln("services stopped: ", srv.Instances.Stopped)
	case "paused":
		srv.Instances.Paused = append(srv.Instances.Paused, id)
	default:
		log.WithFields(log.Fields{
			"error":   "Unknown container state: " + status,
			"service": srvName,
			"id":      id,
		}).Warnln("Cannot add resource to monitor")
	}

	log.WithFields(log.Fields{
		"status":  status,
		"service": srvName,
	}).Infoln("Added resource to monitor")
}

func findIdIndex(id string, instances []string) (int, error) {
	for index, v := range instances {
		if v == id {
			return index, nil
		}
	}

	return -1, ErrNoIndexById
}

func removeResource(id string, stats *GruStats, hist *statsHistory) {
	srv, err := service.GetServiceById(id)
	if err != nil {
		log.Warningln("Cannor remove resource: service unknown")
		return
	}

	running := srv.Instances.Running
	pending := srv.Instances.Pending

	index, err := findIdIndex(id, running)
	if err != nil {
		// If it is not runnig it should be pending
		index, err = findIdIndex(id, pending)
		if err != nil {
			log.WithField("id", id).Errorln("Cannot find pending container to stop")
			return
		}
		pending = append(pending[:index], pending[index+1:]...)
		srv.Instances.Pending = pending
	} else {
		running = append(running[:index], running[index+1:]...)
		srv.Instances.Running = running
	}

	srv.Instances.Stopped = append(srv.Instances.Stopped, id)

	// Upating Event stats
	srvStats := stats.Service[srv.Name]
	srvStats.Events.Stop = append(srvStats.Events.Stop, id)
	stats.Service[srv.Name] = srvStats

	delete(stats.Instance, id)
	delete(hist.instance, id)

	res.FreeInstanceCores(id)

	log.WithFields(log.Fields{
		"service": srv.Name,
		"id":      id,
	}).Infoln("Removed instance")
}

func startMonitorLog(id string) {
	var optionsLog = dockerclient.LogOptions{Follow: true, Stdout: true, Stderr: true, Tail: 1}
	contLog, err := container.Docker().Client.ContainerLogs(id, &optionsLog)
	if err != nil {
		log.WithField("error", err).Errorln("Cannot start log monitoring on container ", id)
	} else {
		metric.Manager().StartCollector(contLog)
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

func statCallBack(id string, stats *dockerclient.Stats, ec chan error, args ...interface{}) {
	if instHist, ok := history.instance[id]; ok {
		// Instance stats update

		// Cpu history usage update
		totCpu := float64(stats.CpuStats.CpuUsage.TotalUsage)
		sysCpu := float64(stats.CpuStats.SystemUsage)
		instHist.cpu.totalUsage.PushBack(totCpu)
		instHist.cpu.sysUsage.PushBack(sysCpu)

		// Memory usage update
		mem := float64(stats.MemoryStats.Usage)
		instHist.mem.PushBack(mem)

		history.instance[id] = instHist
	} else {
		log.WithField("id", id).Debugln("Cannot find history of instance")
	}

}

func monitorError(err error) {
	log.WithField("err", err).Errorln("Error monitoring containers")
	c_err <- err
}

func Stop() {
	monitorActive = false
	log.Warnln("Autonomic monitor stopped")
	c_stop <- struct{}{}
}
