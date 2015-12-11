package service

import (
	"errors"
	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

const c_INFINITE = 1000000000

var (
	services         []Service
	ErrNoSuchService = errors.New("Service does not exists")
)

// func LoadServices(path string) error {
// 	folder, err := ioutil.ReadDir(path)
// 	if err != nil {
// 		log.Errorln("Error opening services folder", err.Error())
// 		return err
// 	}

// 	for _, file := range folder {
// 		var service Service
// 		filep := path + string(filepath.Separator) + file.Name()
// 		log.Debugln("reading file ", filep)
// 		tmp, _ := ioutil.ReadFile(filep)
// 		err = json.Unmarshal(tmp, &service)
// 		if err != nil {
// 			log.WithFields(log.Fields{
// 				"file": file.Name(),
// 				"err":  err,
// 			}).Errorln("Error unmarshaling service file")
// 		} else {
// 			checkService(&service)
// 			services = append(services, service)
// 		}
// 	}

// 	return nil
// }

func Initialize(list []Service) {
	for i := 0; i < len(list); i++ {
		checkService(&list[i])
	}
	services = list
}

func checkService(service *Service) {
	name := service.Name
	checkConstraints(name, &service.Constraints)
	checkConfiguration(name, &service.Configuration)
}

func checkConstraints(name string, cnstrnts *Constraints) {
	if cnstrnts.MaxRespTime == 0 {
		log.WithField("service", name).Warnln("Max Response Time is 0. Setting it to 'infinite'")
		cnstrnts.MaxRespTime = c_INFINITE
	}
}

func checkConfiguration(name string, conf *Config) {
	if conf.Memory == "" {
		log.WithField("service", name).Warnln("Memory limit not set. Service will use all the memory available")
	}

	if conf.CpusetCpus == "" {
		log.WithField("service", name).Warnln("Cores not assigned. Service will use all the cores")
	}
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
