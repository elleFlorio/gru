package monitor

import (
	"strings"
	"time"

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

// Manage channels using the proper package
var ch_mnt_stats_err chan error
var ch_mnt_events_err chan error

func init() {
	ch_mnt_stats_err = make(chan error)
	ch_mnt_events_err = make(chan error)
}

func Start(cError chan error, cStop chan struct{}) {
	go startMonitoring(cError, cStop)
}

func startMonitoring(cError chan error, cStop chan struct{}) {
	log.Debugln("Running autonomic monitoring")
	metric.Manager().Start()

	monitorActive = true
	c_err = cError
	cStop = cStop

	container.Docker().Client.StartMonitorEvents(eventCallback, ch_mnt_events_err)

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
			setServiceInstanceResources(srv.Name, c.Id)
			addInstance(c.Id, srv.Name, status, &gruStats, &history)
			container.Docker().Client.StartMonitorStats(c.Id, statCallBack, ch_mnt_stats_err)
			if status == "pending" {
				startMonitorLog(c.Id)
			}
		}
	}

	for monitorActive {
		select {
		case err := <-ch_mnt_events_err:
			log.WithField("err", err).Debugln("Error monitoring containers events")
			c_err <- err
		case err := <-ch_mnt_stats_err:
			log.WithField("err", err).Debugln("Error monitoring containers stats")
			c_err <- err
		}
	}
}

// Events are: attach, commit, copy, create, destroy, die, exec_create, exec_start, export, kill, oom, pause, rename, resize, restart, start, stop, top, unpause, update
func eventCallback(event *dockerclient.Event, ec chan error, args ...interface{}) {
	log.Debugln("Received event")
	// By now we do not handle events with type != container
	if event.Type != "container" {
		log.WithField("type", event.Type).Debugln("Received event with type different from 'container'")
		return
	}

	srv, err := service.GetServiceByImage(event.From)
	if err != nil {
		log.WithFields(log.Fields{
			"err":   err,
			"event": event,
		}).Warnln("Cannot handle event")
		return
	}

	switch event.Status {
	case "create":
		log.WithField("image", event.From).Debugln("Received create signal")
		setServiceInstanceResources(srv.Name, event.ID)
		container.Docker().Client.StartMonitorStats(event.ID, statCallBack, ch_mnt_stats_err)
	case "start":
		log.WithField("image", event.From).Debugln("Received start signal")
		addInstance(event.ID, srv.Name, "pending", &gruStats, &history)
		startMonitorLog(event.ID)
	case "stop":
		log.WithField("image", event.From).Debugln("Received stop signal")
	case "kill":
		log.WithField("image", event.From).Debugln("Received kill signal")
	case "die":
		log.WithField("image", event.From).Debugln("Received die signal")
		stopInstance(event.ID, &gruStats, &history)
	case "destroy":
		log.WithField("id", event.ID).Debugln("Received destroy signal")
		freeServiceInstanceResources(srv.Name, event.ID)
		removeInstance(event.ID, &gruStats, &history)
	default:
		log.WithFields(log.Fields{
			"err":   "event not handled",
			"event": event.Status,
			"image": event.From,
		}).Debugln("Received unknown signal")
	}

}

func setServiceInstanceResources(name string, id string) {
	var err error

	log.Debugln("Setting new instance resources")
	// This is needed otherwise dockerclient does not
	// return the correct container information
	time.Sleep(100 * time.Millisecond)

	info, err := container.Docker().Client.InspectContainer(id)
	if err != nil {
		log.WithFields(log.Fields{
			"id":  id,
			"err": err,
		}).Errorln("Error setting instance resources")
	}

	cpusetcpus := info.HostConfig.CpusetCpus
	portBindings := createPortBindings(info.HostConfig.PortBindings)

	log.WithFields(log.Fields{
		"service":      name,
		"cpusetcpus":   cpusetcpus,
		"portbindings": portBindings,
	}).Debugln("New instance respources")

	err = res.CheckAndSetSpecificCores(cpusetcpus, id)
	if err != nil {
		log.WithFields(log.Fields{
			"service": name,
			"id":      id,
			"cpus":    cpusetcpus,
			"err":     err,
		}).Errorln("Error assigning CPU resources to new instance")
	}

	err = res.AssignSpecifiPortsToService(name, id, portBindings)
	if err != nil {
		log.WithFields(log.Fields{
			"service":  name,
			"id":       id,
			"bindings": portBindings,
			"err":      err,
		}).Errorln("Error assigning port resources to new instance")
	}
}

func addInstance(id string, srvName string, status string, stats *GruStats, hist *statsHistory) {
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

		service.RegisterServiceInstanceId(srvName, id)
		service.KeepAlive(srvName, id)

	case "stopped":
		srv.Instances.Stopped = append(srv.Instances.Stopped, id)
		service.UnregisterServiceInstance(srvName, id)
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

func stopInstance(id string, stats *GruStats, hist *statsHistory) {
	srv, err := service.GetServiceById(id)
	if err != nil {
		log.Warningln("Cannor stop instance: service unknown")
		return
	}

	running := srv.Instances.Running
	pending := srv.Instances.Pending

	index, err := findIdIndex(id, running)
	if err != nil {
		// If it is not runnig it should be pending
		index, err = findIdIndex(id, pending)
		if err != nil {
			log.WithField("id", id).Debugln("Cannot find pending container to stop")
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

	service.UnregisterServiceInstance(srv.Name, id)

	log.WithFields(log.Fields{
		"service": srv.Name,
		"id":      id,
	}).Infoln("stopped instance")
}

func freeServiceInstanceResources(name string, id string) {
	res.FreeInstanceCores(id)
	res.FreePortsFromService(name, id)
}

func removeInstance(id string, stats *GruStats, hist *statsHistory) {
	srv, err := service.GetServiceById(id)
	if err != nil {
		log.Warnln("Cannor remove instance: service unknown")
		return
	}

	stopInstance(id, stats, hist)
	stopped := srv.Instances.Stopped
	index, err := findIdIndex(id, stopped)
	if err != nil {
		log.Warnln("Cannot find stopped container to remove")
		return
	}
	stopped = append(stopped[:index], stopped[index+1:]...)
	srv.Instances.Stopped = stopped

	service.RemoveInstanceAddress(id)
}

func findIdIndex(id string, instances []string) (int, error) {
	for index, v := range instances {
		if v == id {
			return index, nil
		}
	}

	return -1, ErrNoIndexById
}

func createPortBindings(dockerBindings map[string][]dockerclient.PortBinding) map[string][]string {
	portBindings := make(map[string][]string)

	for guestTcp, bindings := range dockerBindings {
		guest := strings.Split(guestTcp, "/")[0]
		hosts := make([]string, 0, len(bindings))
		for _, host := range bindings {
			hosts = append(hosts, host.HostPort)
		}
		portBindings[guest] = hosts
	}

	return portBindings
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
	log.WithField("err", err).Debugln("Error monitoring containers")
	c_err <- err
}

func Stop() {
	monitorActive = false
	log.Warnln("Autonomic monitor stopped")
	c_stop <- struct{}{}
}
