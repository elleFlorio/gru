package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

type Service struct {
	Name        string
	Type        string
	Image       string
	CpuAvg      float64
	Instances   map[string]Instance
	Constraints Constraints
}

type Instance struct {
	Id  string
	Cpu uint64
}

type Constraints struct {
	CpuMax    float64
	CpuMin    float64
	MinActive int
	MaxActive int
}

var (
	services         []Service
	ErrNoSuchService = errors.New("Service not exists")
)

func LoadServices(path string) ([]Service, error) {
	folder, err := ioutil.ReadDir(path)
	if err != nil {
		log.Errorln("Error opening services folder", err.Error())
		return nil, err
	}

	for _, file := range folder {
		filep := path + string(filepath.Separator) + file.Name()
		log.Debugln("reading file ", filep)
		service := Service{Instances: make(map[string]Instance)}
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

func GetServiceByName(sName string) (*Service, error) {

	for _, service := range services {
		if service.Name == sName {
			return &service, nil
		}
	}

	return nil, ErrNoSuchService
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

func GetServiceByImage(sImg string) (*Service, error) {

	for _, service := range services {
		if service.Image == sImg {
			return &service, nil
		}
	}

	return nil, ErrNoSuchService
}

func GetServiceByInstanceId(sId string) (*Service, error) {

	for _, service := range services {
		if _, ok := service.Instances[sId]; ok {
			return &service, nil
		}
	}

	return nil, ErrNoSuchService
}

func CleanServices() {
	services = make([]Service, 0)
}
