package service

import (
	"errors"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	cfg "github.com/elleFlorio/gru/configuration"
)

const c_INFINITE = 1000000000

var (
	ErrNoSuchService = errors.New("Service does not exists")
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
	return getServiceBy("Name", sName)
}

func GetServiceByImage(sImg string) (*cfg.Service, error) {
	return getServiceBy("Image", sImg)
}

func GetServiceById(id string) (*cfg.Service, error) {
	return getServiceBy("id", id)
}

func getServiceBy(field string, value string) (*cfg.Service, error) {
	services := cfg.GetServices()
	for i := 0; i < len(services); i++ {
		switch field {
		case "Name":
			if services[i].Name == value {
				return &services[i], nil
			}
		case "Image":
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
