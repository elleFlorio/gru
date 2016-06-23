package service

import (
	"errors"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/enum"
)

const c_INFINITE = 1000000000

var (
	ErrNoSuchService error = errors.New("Service does not exists")
	ErrNoIndexById   error = errors.New("No index for such Id")
	ErrUnknownStatus error = errors.New("Unknown instance status")
)

func CheckServices(services []cfg.Service) {
	for i := 0; i < len(services); i++ {
		checkService(&services[i])
	}
}

func checkService(service *cfg.Service) {
	name := service.Name
	checkConstraints(name, &service.Constraints)
	checkConfiguration(name, &service.Docker)
}

func checkConstraints(name string, cnstrnts *cfg.ServiceConstraints) {
	if cnstrnts.MaxRespTime == 0 {
		log.WithField("service", name).Warnln("Max Response Time is 0. Setting it to 'infinite'")
		cnstrnts.MaxRespTime = c_INFINITE
	}
}

func checkConfiguration(name string, conf *cfg.ServiceDocker) {
	if conf.Memory == "" {
		log.WithField("service", name).Warnln("Memory limit not set. Service will use all the memory available")
	}

	if conf.CpusetCpus == "" {
		log.WithField("service", name).Warnln("Cores not assigned. Service will use all the cores")
	}
}

func List() []string {
	names := []string{}
	services := cfg.GetServices()

	for _, service := range services {
		names = append(names, service.Name)
	}

	return names
}

func GetServiceByType(sType string) []cfg.Service {
	byType := make([]cfg.Service, 0)
	services := cfg.GetServices()

	for i := 0; i < len(services); i++ {
		if services[i].Type == sType {
			byType = append(byType, services[i])
		}
	}

	return byType
}

func GetServiceByName(sName string) (*cfg.Service, error) {
	return getServiceBy("name", sName)
}

func GetServiceByImage(sImg string) (*cfg.Service, error) {
	return getServiceBy("image", sImg)
}

func GetServiceById(id string) (*cfg.Service, error) {
	return getServiceBy("id", id)
}

func getServiceBy(field string, value string) (*cfg.Service, error) {
	services := cfg.GetServices()
	for i := 0; i < len(services); i++ {
		switch field {
		case "name":
			if services[i].Name == value {
				return &services[i], nil
			}
		case "image":
			if services[i].Image == value {
				return &services[i], nil
			}
		case "id":
			instances := services[i].Instances.All
			for _, instance := range instances {
				if instance == value {
					return &services[i], nil
				}
			}
		}
	}

	return nil, ErrNoSuchService
}

func GetServiceInstanceStatus(name string, instance string) enum.Status {
	service, err := getServiceBy("name", name)
	if err != nil {
		log.WithField("service", name).Errorln("Cannot get service instance status: service unknown")
		return enum.UNKNOWN
	}

	switch {
	case contains(service.Instances.Pending, instance):
		return enum.PENDING
	case contains(service.Instances.Running, instance):
		return enum.RUNNING
	case contains(service.Instances.Stopped, instance):
		return enum.STOPPED
	case contains(service.Instances.Paused, instance):
		return enum.PAUSED
	default:
		return enum.UNKNOWN
	}
}

func contains(list []string, instance string) bool {
	for _, elem := range list {
		if elem == instance {
			return true
		}
	}

	return false
}

func AddServiceInstance(name string, instance string, status enum.Status) error {
	service, err := getServiceBy("name", name)
	if err != nil {
		log.WithField("service", name).Errorln("Cannot add service instance: service unknown")
		return err
	}

	service.Instances.All = append(service.Instances.All, instance)
	switch status {
	case enum.PENDING:
		service.Instances.Pending = append(service.Instances.Pending, instance)
	case enum.STOPPED:
		service.Instances.Stopped = append(service.Instances.Stopped, instance)
	case enum.PAUSED:
		service.Instances.Paused = append(service.Instances.Paused, instance)
	}

	return nil
}

func ChangeServiceInstanceStatus(name string, instance string, prev enum.Status, upd enum.Status) error {
	service, err := getServiceBy("name", name)
	if err != nil {
		log.WithField("service", name).Errorln("Cannot change service instance status: service unknown")
		return err
	}

	if prev == enum.UNKNOWN {
		log.WithField("instance", instance).Errorln("Cannot change service instance status: unknown status")
		return ErrUnknownStatus
	}

	switch prev {
	case enum.PENDING:
		instIndex, err := findIdIndex(instance, service.Instances.Pending)
		if err != nil {
			log.WithField("instance", instance).Errorln("Cannot change service instance status: instance unknown")
			return err
		}
		service.Instances.Pending = append(service.Instances.Pending[:instIndex],
			service.Instances.Pending[instIndex+1:]...)
	case enum.RUNNING:
		instIndex, err := findIdIndex(instance, service.Instances.Running)
		if err != nil {
			log.WithField("instance", instance).Errorln("Cannot change service instance status: instance unknown")
			return err
		}
		service.Instances.Running = append(service.Instances.Running[:instIndex],
			service.Instances.Running[instIndex+1:]...)
	case enum.STOPPED:
		instIndex, err := findIdIndex(instance, service.Instances.Stopped)
		if err != nil {
			log.WithField("instance", instance).Errorln("Cannot change service instance status: instance unknown")
			return err
		}
		service.Instances.Stopped = append(service.Instances.Stopped[:instIndex],
			service.Instances.Stopped[instIndex+1:]...)
	case enum.PAUSED:
		instIndex, err := findIdIndex(instance, service.Instances.Paused)
		if err != nil {
			log.WithField("instance", instance).Errorln("Cannot change service instance status: instance unknown")
			return err
		}
		service.Instances.Paused = append(service.Instances.Paused[:instIndex],
			service.Instances.Paused[instIndex+1:]...)
	}

	switch upd {
	case enum.PENDING:
		service.Instances.Pending = append(service.Instances.Pending, instance)
	case enum.RUNNING:
		service.Instances.Running = append(service.Instances.Running, instance)
	case enum.STOPPED:
		service.Instances.Stopped = append(service.Instances.Stopped, instance)
	case enum.PAUSED:
		service.Instances.Paused = append(service.Instances.Paused, instance)
	}

	return nil

}

func RemoveServiceInstance(name string, instance string) error {
	service, err := getServiceBy("name", name)
	if err != nil {
		log.WithField("service", name).Errorln("Cannot remove service instance: service unknown")
		return err
	}

	instIndexStop, err := findIdIndex(instance, service.Instances.Stopped)
	if err != nil {
		log.WithField("instance", instance).Errorln("Cannot remove service instance: instance unknown")
		return err
	}
	instIndexAll, err := findIdIndex(instance, service.Instances.All)
	if err != nil {
		log.WithField("instance", instance).Errorln("Cannot remove service instance: instance unknown")
		return err
	}

	service.Instances.Stopped = append(service.Instances.Stopped[:instIndexStop],
		service.Instances.Stopped[instIndexStop+1:]...)
	service.Instances.All = append(service.Instances.All[:instIndexAll],
		service.Instances.All[instIndexAll+1:]...)

	return nil
}

func findIdIndex(id string, instances []string) (int, error) {
	for index, v := range instances {
		if v == id {
			return index, nil
		}
	}

	return -1, ErrNoIndexById
}

func IsServiceActive(service string) bool {
	srv, err := getServiceBy("name", service)
	if err != nil {
		log.WithField("service", service).Errorln("Cannot determine if service is active: unknown service")
		return false
	}

	return (len(srv.Instances.Pending) + len(srv.Instances.Running)) > 0
}

func GetServiceExpressionsList(service string) []string {
	srv, err := getServiceBy("name", service)
	if err != nil {
		log.WithField("service", service).Errorln("Cannot return service expressions: unknown service")
		return []string{}
	}

	return srv.Expressions
}
