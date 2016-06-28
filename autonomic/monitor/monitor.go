package monitor

import (
	"fmt"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/samalba/dockerclient"

	evt "github.com/elleFlorio/gru/autonomic/monitor/event"
	lgr "github.com/elleFlorio/gru/autonomic/monitor/logreader"
	mtr "github.com/elleFlorio/gru/autonomic/monitor/metric"
	chn "github.com/elleFlorio/gru/channels"
	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/container"
	"github.com/elleFlorio/gru/data"
	"github.com/elleFlorio/gru/enum"
	srv "github.com/elleFlorio/gru/service"
	"github.com/elleFlorio/gru/utils"
)

// Add memory
type instanceMetricBuffer struct {
	cpuInst utils.Buffer
	cpuSys  utils.Buffer
}

const c_B_SIZE = 20
const c_MTR_THR = 20

var (
	stats            data.GruStats
	instBuffer       map[string]instanceMetricBuffer
	enableLogReading bool

	ch_mnt_stats_err  chan error
	ch_mnt_events_err chan error
	ch_stop           chan struct{}
)

func init() {
	stats = data.GruStats{}
	instBuffer = make(map[string]instanceMetricBuffer)

	ch_mnt_stats_err = make(chan error)
	ch_mnt_events_err = make(chan error)
	ch_stop = make(chan struct{})
}

func StartMonitor() {
	initiailizeMonitoring()
	go startMonitoring()
}

func StopMonitor() {
	ch_stop <- struct{}{}
	log.Warnln("Autonomic monitor stopped")
}

func Run() data.GruStats {
	log.WithField("status", "init").Debugln("Gru Monitor")
	defer log.WithField("status", "done").Debugln("Gru Monitor")

	services := srv.List()
	updateRunningInstances(services, c_MTR_THR)
	updateSystemInstances(services)
	metrics := mtr.GetMetricsStats()
	events := evt.GetEventsStats()
	stats.Metrics = metrics
	stats.Events = events
	data.SaveStats(stats)
	displayStatsOfServices(stats)
	return stats
}

func initiailizeMonitoring() {
	defer log.Infoln("Initializing autonomic monitoring")
	ch_aut_err := chn.GetAutonomicErrChannel()
	enableLogReading = cfg.GetAgentAutonomic().EnableLogReading
	mtr.Initialize(srv.List())

	// Start log reader if needed
	if enableLogReading {
		lgr.Initialize(srv.List())
		lgr.StartLogReader()
		log.WithField("logreader", enableLogReading).Debugln("Log reading is enabled")
	}

	// Get the list of containers (running or not) to monitor
	containers, err := container.Docker().Client.ListContainers(true, false, "")
	if err != nil {
		log.WithField("err", err).Debugln("Error monitoring containers")
		ch_aut_err <- err
	}

	// Start the monitor for each configured service
	for _, c := range containers {
		info, _ := container.Docker().Client.InspectContainer(c.Id)
		status := getContainerStatus(info)
		service, err := srv.GetServiceByImage(c.Image)
		if err != nil {
			log.WithFields(log.Fields{
				"err":   err,
				"image": c.Image,
			}).Warningln("Error monitoring service")
		} else {
			e := evt.Event{
				Service:  service.Name,
				Image:    c.Image,
				Instance: c.Id,
				Status:   status,
			}

			evt.HandleCreateEvent(e)
			evt.HanldeStartEvent(e)
			mtr.AddInstance(c.Id)
			if _, ok := instBuffer[c.Id]; !ok {
				instBuffer[c.Id] = instanceMetricBuffer{
					cpuInst: utils.BuildBuffer(c_B_SIZE),
					cpuSys:  utils.BuildBuffer(c_B_SIZE),
				}
			}
			container.Docker().Client.StartMonitorStats(c.Id, statCallBack, ch_mnt_stats_err)
			if status == enum.PENDING && enableLogReading {
				startMonitorLog(c.Id)
			}
		}
	}
}

func startMonitoring() {
	log.Infoln("Running autonomic monitoring")
	ch_aut_err := chn.GetAutonomicErrChannel()

	container.Docker().Client.StartMonitorEvents(eventCallback, ch_mnt_events_err)
	for {
		select {
		case err := <-ch_mnt_events_err:
			log.WithField("err", err).Fatalln("Error monitoring containers events")
			ch_aut_err <- err
		case err := <-ch_mnt_stats_err:
			log.WithField("err", err).Debugln("Error monitoring containers stats")
			ch_aut_err <- err
		case <-ch_stop:
			return
		}
	}
}

func getContainerStatus(info *dockerclient.ContainerInfo) enum.Status {
	switch {
	case info.State.Running:
		return enum.PENDING
	case info.State.Paused:
		return enum.PAUSED
	}

	return enum.STOPPED
}

// Events are: attach, commit, copy, create, destroy, die, exec_create, exec_start, export, kill, oom, pause, rename, resize, restart, start, stop, top, unpause, update
func eventCallback(event *dockerclient.Event, ec chan error, args ...interface{}) {
	log.Debugln("Received event")
	// By now we do not handle events with type != container
	if event.Type != "container" {
		log.WithField("type", event.Type).Debugln("Received event with type different from 'container'")
		return
	}

	service, err := srv.GetServiceByImage(event.From)
	if err != nil {
		log.WithFields(log.Fields{
			"err":   err,
			"event": event,
		}).Warnln("Cannot handle event")
		return
	}

	e := evt.Event{
		Service:  service.Name,
		Image:    event.From,
		Instance: event.ID,
		Type:     event.Type,
	}

	switch event.Status {
	case "create":
		log.WithField("image", e.Image).Debugln("Received create signal")
		evt.HandleCreateEvent(e)
		container.Docker().Client.StartMonitorStats(e.Instance, statCallBack, ch_mnt_stats_err)
	case "start":
		log.WithField("image", e.Image).Debugln("Received start signal")
		if _, ok := instBuffer[e.Instance]; !ok {
			instBuffer[e.Instance] = instanceMetricBuffer{
				cpuInst: utils.BuildBuffer(c_B_SIZE),
				cpuSys:  utils.BuildBuffer(c_B_SIZE),
			}
		}
		e.Status = enum.PENDING
		evt.HanldeStartEvent(e)
		mtr.AddInstance(e.Instance)
		if enableLogReading {
			startMonitorLog(event.ID)
		}
	case "stop":
		log.WithField("image", e.Image).Debugln("Received stop signal")
	case "kill":
		log.WithField("image", e.Image).Debugln("Received kill signal")
	case "die":
		log.WithField("image", e.Image).Debugln("Received die signal")
		delete(instBuffer, e.Instance)
		mtr.RemoveInstance(e.Instance)
		evt.HandleStopEvent(e)
	case "destroy":
		log.WithField("id", e.Instance).Debugln("Received destroy signal")
		evt.HandleRemoveEvent(e)
	default:
		log.WithFields(log.Fields{
			"err":   "event not handled",
			"event": event.Status,
			"image": event.From,
		}).Debugln("Received unknown signal")
	}

	log.Debugln("Event handled")

}

func startMonitorLog(id string) {
	var optionsLog = dockerclient.LogOptions{Follow: true, Stdout: true, Stderr: true, Tail: 1}
	contLog, err := container.Docker().Client.ContainerLogs(id, &optionsLog)
	if err != nil {
		log.WithField("error", err).Errorln("Cannot start log monitoring on container ", id)
	} else {
		lgr.StartCollector(contLog)
	}
}

func statCallBack(id string, stats *dockerclient.Stats, ec chan error, args ...interface{}) {

	metricBuffer := instBuffer[id]
	cpuInst := float64(stats.CpuStats.CpuUsage.TotalUsage)
	cpuSys := float64(stats.CpuStats.SystemUsage)

	toAddInst := metricBuffer.cpuInst.PushValue(cpuInst)
	toAddSys := metricBuffer.cpuSys.PushValue(cpuSys)

	instBuffer[id] = metricBuffer

	if toAddInst != nil && toAddSys != nil {
		mtr.UpdateCpuMetric(id, toAddInst, toAddSys)
	}

	// TODO - ADD MEMORY

}

func updateRunningInstances(services []string, threshold int) {
	for _, name := range services {
		service, _ := srv.GetServiceByName(name)
		pending := service.Instances.Pending

		for _, instance := range pending {
			if mtr.IsReadyForRunning(instance, threshold) {
				// TODO
				e := evt.Event{
					Service:  name,
					Instance: instance,
					Status:   enum.PENDING,
				}
				evt.HandlePromoteEvent(e)
				log.WithFields(log.Fields{
					"service":  name,
					"instance": instance,
				}).Debugln("Promoted resource to running state")
			}
		}
	}
}

func updateSystemInstances(services []string) {
	cfg.ClearNodeInstances()
	instances := cfg.GetNodeInstances()
	for _, name := range services {
		service, _ := srv.GetServiceByName(name)
		instances.All = append(instances.All, service.Instances.All...)
		instances.Pending = append(instances.Pending, service.Instances.Pending...)
		instances.Running = append(instances.Running, service.Instances.Running...)
		instances.Stopped = append(instances.Stopped, service.Instances.Stopped...)
		instances.Paused = append(instances.Paused, service.Instances.Paused...)
	}
}

func displayStatsOfServices(stats data.GruStats) {
	for name, value := range stats.Metrics.Service {
		service, _ := srv.GetServiceByName(name)
		log.WithFields(log.Fields{
			"pending:": len(service.Instances.Pending),
			"running:": len(service.Instances.Running),
			"stopped:": len(service.Instances.Stopped),
			"paused:":  len(service.Instances.Paused),
			"cpu avg":  fmt.Sprintf("%.2f", value.BaseMetrics[enum.METRIC_CPU_AVG.ToString()]),
			"mem avg":  fmt.Sprintf("%.2f", value.BaseMetrics[enum.METRIC_MEM_AVG.ToString()]),
		}).Infoln("Stats computed: ", name)
	}
}
