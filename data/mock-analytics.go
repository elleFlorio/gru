package data

import (
	cfg "github.com/elleFlorio/gru/configuration"
)

func CreateMockAnalytics() GruAnalytics {
	mockAnalytics := GruAnalytics{
		Service: make(map[string]ServiceAnalytics),
	}

	service1 := "service1"
	service2 := "service2"

	cpu1 := 0.4
	mem1 := 0.1
	avail1 := 1.0
	res1 := ResourcesAnalytics{
		Cpu:       cpu1,
		Memory:    mem1,
		Available: avail1,
	}

	load1 := 0.3
	status1 := cfg.ServiceStatus{
		Running: []string{"service1_1"},
		Stopped: []string{"service1_2"},
	}
	health1 := 0.9

	analytics1 := ServiceAnalytics{
		Load:      load1,
		Resources: res1,
		Instances: status1,
		Health:    health1,
	}

	cpu2 := 0.8
	mem2 := 0.4
	avail2 := 0.0
	res2 := ResourcesAnalytics{
		Cpu:       cpu2,
		Memory:    mem2,
		Available: avail2,
	}

	load2 := 0.8
	status2 := cfg.ServiceStatus{
		Running: []string{"service2_1", "service2_2"},
	}

	health2 := 0.6

	analytics2 := ServiceAnalytics{
		Load:      load2,
		Resources: res2,
		Instances: status2,
		Health:    health2,
	}

	mockAnalytics.Service[service1] = analytics1
	mockAnalytics.Service[service2] = analytics2

	systService := []string{service1, service2}
	systCpu := 0.8
	systMem := 0.5
	systRes := ResourcesAnalytics{
		Cpu:    systCpu,
		Memory: systMem,
	}
	systInstances := cfg.ServiceStatus{
		Running: []string{"service1_1", "service2_1", "service2_2"},
		Stopped: []string{"service1_2"},
	}
	systHealth := 0.8

	systAnalytics := SystemAnalytics{
		Services:  systService,
		Resources: systRes,
		Instances: systInstances,
		Health:    systHealth,
	}

	mockAnalytics.System = systAnalytics

	clustAnalytics := ClusterAnalytics{
		Services:  systService,
		Resources: systRes,
		Health:    systHealth,
	}

	mockAnalytics.Cluster = clustAnalytics

	return mockAnalytics
}

func StoreMockAnalytics() {
	SaveAnalytics(CreateMockAnalytics())
}
