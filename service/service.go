package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

var (
	services         []Service
	ErrNoSuchService = errors.New("Service does not exists")
)

func LoadServices(path string) error {
	folder, err := ioutil.ReadDir(path)
	if err != nil {
		log.Errorln("Error opening services folder", err.Error())
		return err
	}

	for _, file := range folder {
		var service Service
		filep := path + string(filepath.Separator) + file.Name()
		log.Debugln("reading file ", filep)
		tmp, _ := ioutil.ReadFile(filep)
		err = json.Unmarshal(tmp, &service)
		if err != nil {
			log.WithFields(log.Fields{
				"file":  file.Name(),
				"error": err,
			}).Errorln("Error unmarshaling service file")
		} else {
			services = append(services, service)
		}
	}

	log.Infoln("Services loading complete. Loaded files: ", len(services))

	return nil
}

func List() []string {
	names := []string{}

	for _, service := range services {
		names = append(names, service.Name)
	}

	return names
}

func GetServices() []Service {
	return services
}

func GetServiceByType(sType string) []Service {
	byType := make([]Service, 0)

	for _, service := range services {
		if service.Type == sType {
			byType = append(byType, service)
		}
	}

	return byType
}

func GetServiceByName(sName string) (*Service, error) {
	return getServiceBy("Name", sName)
}

func GetServiceByImage(sImg string) (*Service, error) {
	return getServiceBy("Image", sImg)
}

func getServiceBy(field string, value string) (*Service, error) {
	for _, service := range services {
		switch field {
		case "Name":
			if service.Name == value {
				return &service, nil
			}
		case "Image":
			if service.Image == value {
				return &service, nil
			}
		}
	}

	return nil, ErrNoSuchService
}

func AddServices(newServices []Service) {
	services = append(services, newServices...)
}

func RemoveServices(rmServices []string) {
	indexes := make([]int, len(rmServices), len(rmServices))

	for i, rmService := range rmServices {
		for j, service := range services {
			if service.Name == rmService {
				indexes[i] = j
			}
		}
	}

	for _, index := range indexes {
		services = append(services[:index], services[index+1:]...)
	}
}

func UpdateServices(newServices []Service) {
	CleanServices()
	services = newServices
}

func CleanServices() {
	services = make([]Service, 0)
}
