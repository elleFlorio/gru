package service

import (
	cfg "github.com/elleFlorio/gru/configuration"
)

func CreateMockServices() []cfg.Service {
	all1 := []string{"instance1_0, instance1_1", "instance1_2", "instance1_3", "instance1_4"}
	running1 := []string{"instance1_1", "instance1_2"}
	pending1 := []string{"instance1_3", "instance1_4"}
	stopped1 := []string{"instance1_0"}
	paused1 := []string{}
	instances1 := cfg.ServiceStatus{
		All:     all1,
		Running: running1,
		Pending: pending1,
		Stopped: stopped1,
		Paused:  paused1,
	}
	service1 := cfg.Service{
		Name:      "service1",
		Type:      "webserver",
		Image:     "test/tomcat",
		Instances: instances1,
		Constraints: cfg.ServiceConstraints{
			MaxRespTime: 2000,
		},
	}

	all2 := []string{"instance2_1"}
	running2 := []string{"instance2_1"}
	instances2 := cfg.ServiceStatus{
		All:     all2,
		Running: running2,
	}
	service2 := cfg.Service{
		Name:      "service2",
		Type:      "webserver",
		Image:     "test/jetty",
		Instances: instances2,
		Constraints: cfg.ServiceConstraints{
			MaxRespTime: 6000,
		},
	}

	service3 := cfg.Service{
		Name:  "service3",
		Type:  "database",
		Image: "test/mysql",
		Constraints: cfg.ServiceConstraints{
			MaxRespTime: 1000,
		},
	}

	return []cfg.Service{service1, service2, service3}
}
