package event

import (
	"errors"
	"math"
	"sync"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	chn "github.com/elleFlorio/gru/channels"
	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/data"
	"github.com/elleFlorio/gru/enum"
	res "github.com/elleFlorio/gru/resources"
	srv "github.com/elleFlorio/gru/service"
)

const c_MILLISEC = 1000

var (
	events    data.EventStats
	evt_mutex sync.RWMutex

	ErrNoIndexById error = errors.New("No index for such Id")
)

func init() {
	events.Service = make(map[string]data.EventData)
	evt_mutex = sync.RWMutex{}
}

func Initialize(services []string) {
	for _, service := range services {
		events.Service[service] = data.EventData{
			Start: []string{},
			Stop:  []string{},
		}
	}
}

func HandleCreateEvent(e Event) {
	res.SetServiceInstanceResources(e.Service, e.Instance)
}

func HanldeStartEvent(e Event) {
	addInstance(e.Service, e.Instance, e.Status)
	updateAutoLoopTimeInterval()
}

func HandlePromoteEvent(e Event) {
	promotePendingToRunning(e.Service, e.Instance)
}

func HandleStopEvent(e Event) {
	stopInstance(e.Service, e.Instance)
	updateAutoLoopTimeInterval()
}

func HandleRemoveEvent(e Event) {
	freeServiceInstanceResources(e.Service, e.Instance)
	removeInstance(e.Service, e.Instance)
	notifyRemoval()
}

func GetEventsStats() data.EventStats {
	defer clearEvents()
	events_cpy := data.EventStats{
		Service: make(map[string]data.EventData, len(events.Service)),
	}
	for service, values := range events.Service {
		events_cpy.Service[service] = values
	}

	return events_cpy
}

func clearEvents() {
	evt_mutex.Lock()
	for service, eventData := range events.Service {
		eventData.Start = eventData.Start[:0]
		eventData.Stop = eventData.Stop[:0]
		events.Service[service] = eventData
	}
	evt_mutex.Unlock()
}

func addInstance(name string, instance string, status enum.Status) {
	err := srv.AddServiceInstance(name, instance, status)
	if err != nil {
		log.WithField("err", err).Errorln("Cannot add new service instance")
		return
	}

	if status == enum.PENDING {
		evt_mutex.Lock()
		srvEvents := events.Service[name]
		srvEvents.Start = append(srvEvents.Start, instance)
		events.Service[name] = srvEvents
		evt_mutex.Unlock()

		srv.RegisterServiceInstanceId(name, instance)
		srv.KeepAlive(name, instance)
	}

	log.WithFields(log.Fields{
		"status":  status,
		"service": name,
	}).Infoln("Added resource to monitor")
}

func promotePendingToRunning(name string, instance string) {
	srv.ChangeServiceInstanceStatus(name, instance, enum.PENDING, enum.RUNNING)
}

func stopInstance(name string, instance string) {
	log.Debugln("Changing instance status to stop...")
	status := srv.GetServiceInstanceStatus(name, instance)
	if status == enum.STOPPED {
		log.Debugln("Instance already stopped")
		return
	}

	err := srv.ChangeServiceInstanceStatus(name, instance, status, enum.STOPPED)
	if err != nil {
		log.WithField("err", err).Errorln("Cannot stop service instance")
		return
	}

	log.Debugln("Updating events...")
	evt_mutex.Lock()
	srvEvents := events.Service[name]
	srvEvents.Stop = append(srvEvents.Stop, instance)
	events.Service[name] = srvEvents
	evt_mutex.Unlock()

	log.Debugln("Unregistering instance...")
	srv.UnregisterServiceInstance(name, instance)

	log.WithFields(log.Fields{
		"service": name,
		"id":      instance,
	}).Infoln("stopped instance")
}

func freeServiceInstanceResources(name string, id string) {
	res.FreeInstanceCores(id)
	res.FreePortsFromService(name, id)
}

func removeInstance(name string, instance string) {
	stopInstance(name, instance)
	err := srv.RemoveServiceInstance(name, instance)
	if err != nil {
		log.WithField("err", err).Errorln("Cannot remove service instance")
		return
	}

	srv.RemoveInstanceAddress(instance)
}

func notifyRemoval() {
	log.Debugln("Checking removal notification...")
	if chn.NeedsRemovalNotification() {
		log.Debugln("Needs removal notification. Sending message...")
		chn.GetRemovalChannel() <- struct{}{}
		log.Debugln("Removal notified")
		chn.SetRemovalNotification(false)
	}
}

func updateAutoLoopTimeInterval() {
	if cfg.GetAgentAutonomic().EnableDynamicLoop {
		interval := 0.0
		for _, service := range srv.GetActiveServices() {
			if mrp, ok := service.Constraints["MAX_RESP_TIME"]; ok {
				if mrp > interval {
					interval = mrp
				}
			} else {
				log.WithField("service", service.Name).Warnln("No MAX_RESP_TIME constraint")
			}
		}

		if interval > 0.0 {
			cfg.GetAgentAutonomic().LoopTimeInterval = int(math.Ceil(interval / c_MILLISEC))
		}
	}
}
