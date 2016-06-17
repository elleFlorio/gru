package event

import (
	"errors"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	chn "github.com/elleFlorio/gru/channels"
	"github.com/elleFlorio/gru/data"
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
	res.SetServiceInstanceResources(e.Service, e.Isntance)
}

func HanldeStartEvent(e Event) {
	addInstance(e.Service, e.Isntance, e.Status)
}

func HandleStopEvent(e Event) {
	stopInstance(e.Isntance)
}

func HandleRemoveEvent(e Event) {
	freeServiceInstanceResources(e.Service, e.Isntance)
	removeInstance(e.Isntance)
	notifyRemoval()
}

func addInstance(name string, instance string, status string) {
	service, _ := srv.GetServiceByName(name)
	_, err := findIdIndex(instance, service.Instances.All)
	if err != nil {
		service.Instances.All = append(service.Instances.All, instance)
	}

	switch status {
	case "running":
		index, err := findIdIndex(instance, service.Instances.Pending)
		service.Instances.Running = append(service.Instances.Running, instance)
		if err != nil {
			log.WithField("error", err).Errorln("Cannot find pending instance to promote running")
		} else {
			service.Instances.Pending = append(
				service.Instances.Pending[:index],
				service.Instances.Pending[index+1:]...)
		}
	case "pending":
		service.Instances.Pending = append(service.Instances.Pending, instance)
		index, err := findIdIndex(instance, service.Instances.Stopped)
		if err != nil {
			log.WithField("error", err).Debugln("Cannot find stopped instance to promote pending")
		} else {
			service.Instances.Stopped = append(
				service.Instances.Stopped[:index],
				service.Instances.Stopped[index+1:]...)
		}

		srvEvents := events.Service[service.Name]
		srvEvents.Start = append(srvEvents.Start, instance)
		events.Service[service.Name] = srvEvents

		srv.RegisterServiceInstanceId(service.Name, instance)
		srv.KeepAlive(service.Name, instance)

	case "stopped":
		service.Instances.Stopped = append(service.Instances.Stopped, instance)
		log.Debugln("services stopped: ", service.Instances.Stopped)
	case "paused":
		service.Instances.Paused = append(service.Instances.Paused, instance)
	default:
		log.WithFields(log.Fields{
			"error":   "Unknown container state: " + status,
			"service": service.Name,
			"id":      instance,
		}).Warnln("Cannot add resource to monitor")
	}

	log.WithFields(log.Fields{
		"status":  status,
		"service": service.Name,
	}).Infoln("Added resource to monitor")
}

func stopInstance(instance string) {
	service, err := srv.GetServiceById(instance)
	if err != nil {
		log.Warningln("Cannor stop instance: service unknown")
		return
	}

	running := service.Instances.Running
	pending := service.Instances.Pending

	index, err := findIdIndex(instance, running)
	if err != nil {
		// If it is not runnig it should be pending
		index, err = findIdIndex(instance, pending)
		if err != nil {
			log.WithField("id", instance).Debugln("Cannot find pending container to stop")
			return
		}
		pending = append(pending[:index], pending[index+1:]...)
		service.Instances.Pending = pending
	} else {
		running = append(running[:index], running[index+1:]...)
		service.Instances.Running = running
	}

	service.Instances.Stopped = append(service.Instances.Stopped, instance)

	srvEvents := events.Service[service.Name]
	srvEvents.Stop = append(srvEvents.Stop, instance)
	events.Service[service.Name] = srvEvents

	srv.UnregisterServiceInstance(service.Name, instance)

	log.WithFields(log.Fields{
		"service": service.Name,
		"id":      instance,
	}).Infoln("stopped instance")
}

func freeServiceInstanceResources(name string, id string) {
	res.FreeInstanceCores(id)
	res.FreePortsFromService(name, id)
}

func removeInstance(instance string) {
	service, err := srv.GetServiceById(instance)
	if err != nil {
		log.WithField("instance", instance).Warnln("Cannot remove instance: service unknown")
		return
	}

	stopInstance(instance)
	stopped := service.Instances.Stopped
	index, err := findIdIndex(instance, stopped)
	if err != nil {
		log.Warnln("Cannot find stopped container to remove")
		return
	}
	stopped = append(stopped[:index], stopped[index+1:]...)
	service.Instances.Stopped = stopped

	srv.RemoveInstanceAddress(instance)
}

func notifyRemoval() {
	if chn.NeedsRemovalNotification() {
		chn.GetRemovalChannel() <- struct{}{}
		chn.SetRemovalNotification(false)
	}
}

func findIdIndex(id string, instances []string) (int, error) {
	for index, v := range instances {
		if v == id {
			return index, nil
		}
	}

	return -1, ErrNoIndexById
}
