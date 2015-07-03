package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/samalba/dockerclient"
)

type Service struct {
	Name            string
	Type            string
	Image           string
	Constraints     Constraints
	ContainerConfig dockerclient.ContainerConfig
	HostConfig      dockerclient.HostConfig
}

type Constraints struct {
	CpuMax    float64
	CpuMin    float64
	MinActive int
	MaxActive int
}

var (
	services         []Service
	ErrNoSuchService = errors.New("Service does not exists")
)

func LoadServices(path string) ([]Service, error) {
	folder, err := ioutil.ReadDir(path)
	if err != nil {
		log.Errorln("Error opening services folder", err.Error())
		return nil, err
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

	return services, nil
}

func List() []string {
	names := []string{}

	for _, service := range services {
		names = append(names, service.Name)
	}

	return names
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

func CleanServices() {
	services = make([]Service, 0)
}
