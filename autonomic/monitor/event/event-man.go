package event

import (
	"errors"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	chn "github.com/elleFlorio/gru/channels"
	"github.com/elleFlorio/gru/data"
	"github.com/elleFlorio/gru/enum"
	res "github.com/elleFlorio/gru/resources"
	srv "github.com/elleFlorio/gru/service"
)

var (
	events data.EventStats

	ErrNoIndexById error = errors.New("No index for such Id")
)

func init() {
	events.Service = make(map[string]data.EventData)
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
}

func HandlePromoteEvent(e Event) {
	promotePendingToRunning(e.Service, e.Instance)
}

func HandleStopEvent(e Event) {
	stopInstance(e.Service, e.Instance)
}

func HandleRemoveEvent(e Event) {
	freeServiceInstanceResources(e.Service, e.Instance)
	removeInstance(e.Service, e.Instance)
	notifyRemoval()
}

// TODO - Need to do a copy
func GetEventsStats() data.EventStats {
	return events
}

func clearEvents() {
	for service, eventData := range events.Service {
		eventData.Start = eventData.Start[:0]
		eventData.Stop = eventData.Stop[:0]
		events.Service[service] = eventData
	}
}

func addInstance(name string, instance string, status enum.Status) {
	err := srv.AddServiceInstance(name, instance, status)
	if err != nil {
		log.WithField("err", err).Errorln("Cannot add new service instance")
		return
	}

	if status == enum.PENDING {
		srvEvents := events.Service[name]
		srvEvents.Start = append(srvEvents.Start, instance)
		events.Service[name] = srvEvents

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
	status := srv.GetServiceInstanceStatus(name, instance)
	err := srv.ChangeServiceInstanceStatus(name, instance, status, enum.STOPPED)
	if err != nil {
		log.WithField("err", err).Errorln("Cannot stop service instance")
		return
	}

	srvEvents := events.Service[name]
	srvEvents.Stop = append(srvEvents.Stop, instance)
	events.Service[name] = srvEvents

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
	if chn.NeedsRemovalNotification() {
		chn.GetRemovalChannel() <- struct{}{}
		chn.SetRemovalNotification(false)
	}
}
