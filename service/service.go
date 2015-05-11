package service

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

type Service struct {
	Name        string
	Type        string
	Image       string
	CpuAvg      float64
	Instances   []Instance
	Constraints Constraints
}

type Instance struct {
	Id  string
	Cpu float64
}

type Constraints struct {
	CpuMax    float64
	CpuMin    float64
	MinActive int
	MaxActive int
}

var services []Service

func LoadServices(path string) ([]Service, error) {
	folder, err := ioutil.ReadDir(path)
	if err != nil {
		log.Errorln("Error opening services folder", err.Error())
		return nil, err
	}

	for _, file := range folder {
		filep := path + string(filepath.Separator) + file.Name()
		log.Debugln("reading file ", filep)
		service := Service{}
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
	names := make([]string, len(services))
	for _, service := range services {
		names = append(names, service.Name)
	}

	return names
}

func GetServiceByName(sName string) []Service {
	byName := make([]Service, 0)

	for _, service := range services {
		if service.Name == sName {
			byName = append(byName, service)
		}
	}

	return byName
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

func GetServiceByImage(sImg string) []Service {
	byImage := make([]Service, 0)

	for _, service := range services {
		if service.Image == sImg {
			byImage = append(byImage, service)
		}
	}

	return byImage
}

func CleanServices() {
	services = make([]Service, 0)
}
