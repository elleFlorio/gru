package service

import (
	cfg "github.com/elleFlorio/gru/configuration"
)

func CreateMockServices() []cfg.Service {
	service1 := cfg.Service{
		Name:  "service1",
		Type:  "webserver",
		Image: "test/tomcat",
		Constraints: cfg.ServiceConstraints{
			MaxRespTime: 2000,
		},
	}

	service2 := cfg.Service{
		Name:  "service2",
		Type:  "webserver",
		Image: "test/jetty",
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
