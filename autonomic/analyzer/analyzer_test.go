package analyzer

import (
	"testing"

	"github.com/elleFlorio/gru/autonomic/monitor"
	"github.com/stretchr/testify/assert"
)

func init() {
	gruAnalytics = GruAnalytics{
		Service:  make(map[string]ServiceAnalytics),
		Instance: make(map[string]InstanceAnalytics),
	}
}

func TestUpdateInstances(t *testing.T) {
	mockStats := createMockStats()
	for name, _ := range mockStats.Service {
		updateInstances(name, &mockStats)
	}

	assert.NotContains(t, gruAnalytics.Service["service1"].Instances.Running, "instance0", "Service 1 should contains instance2")
	inst := []string{}
	for k, _ := range gruAnalytics.Instance {
		inst = append(inst, k)
	}
	assert.NotContains(t, inst, "instance0", "Instances should not contains instance2")
	assert.Contains(t, gruAnalytics.Service["service2"].Instances.Running, "instance3", "Service 2 should contains instance3")

	cleanAnalytics()
}

func TestUpdateAnalytics(t *testing.T) {
	mockStats := createMockStats()
	for name, _ := range mockStats.Service {
		updateAnalytics(name, &mockStats)
	}
	statsCpu := mockStats.Instance["instance1"].Cpu
	analyticsCpu := gruAnalytics.Instance["instance1"].Cpu
	assert.Equal(t, statsCpu, analyticsCpu, "Instance1 stats and analytics should be equal")

	cleanAnalytics()
}

func TestComputeCpuAvg(t *testing.T) {
	mockStats := createMockStats()
	createMockAnalytics()
	for name, _ := range mockStats.Service {
		computeCpuAvg(name, &mockStats)
	}

	assert.Equal(t, 0.25, gruAnalytics.Service["service1"].CpuAvg, "Service1 cpuAvg should be 25%")
	assert.Equal(t, 0.3, gruAnalytics.Service["service2"].CpuAvg, "Service2 cpuAvg should be 30%")
}

func createMockStats() monitor.GruStats {
	all1 := []string{"instance0", "instance1", "instance2", "instance4"}
	running1 := []string{"instance1", "instance2", "instance4"}
	stopped1 := []string{"instance0"}
	events1 := monitor.EventStats{
		Start: []string{"instance4"},
		Stop:  []string{"instance0"},
	}
	instances1 := monitor.InstanceStatus{
		All:     all1,
		Running: running1,
		Stopped: stopped1,
	}
	service1 := monitor.ServiceStats{
		Instances: instances1,
		Events:    events1,
	}

	all2 := []string{"instance3"}
	running2 := []string{"instance3"}
	instances2 := monitor.InstanceStatus{
		All:     all2,
		Running: running2,
	}
	service2 := monitor.ServiceStats{Instances: instances2}
	services := map[string]monitor.ServiceStats{
		"service1": service1,
		"service2": service2,
	}

	instStat1 := monitor.InstanceStats{20000}
	instStat2 := monitor.InstanceStats{60000}
	instStat3 := monitor.InstanceStats{60000}
	instStat4 := monitor.InstanceStats{10000}
	instances := map[string]monitor.InstanceStats{
		"instance1": instStat1,
		"instance2": instStat2,
		"instance3": instStat3,
		"instance4": instStat4,
	}

	system := monitor.SystemStats{15000000}

	mockStats := monitor.GruStats{
		Service:  services,
		Instance: instances,
		System:   system,
	}

	return mockStats
}

func cleanAnalytics() {
	gruAnalytics = GruAnalytics{}
	gruAnalytics.Service = make(map[string]ServiceAnalytics)
	gruAnalytics.Instance = make(map[string]InstanceAnalytics)
}

func createMockAnalytics() {
	all1 := []string{"instance0", "instance1", "instance2", "instance4"}
	running1 := []string{"instance4", "instance1", "instance2"}
	instances1 := InstanceStatus{
		All:     all1,
		Running: running1,
	}
	service1 := ServiceAnalytics{
		CpuAvg:    0.2,
		Instances: instances1}

	all2 := []string{"instance3"}
	running2 := []string{"instance3"}
	instances2 := InstanceStatus{
		All:     all2,
		Running: running2,
	}
	service2 := ServiceAnalytics{
		CpuAvg:    0.4,
		Instances: instances2,
	}
	services := map[string]ServiceAnalytics{
		"service1": service1,
		"service2": service2,
	}

	instAnalytics1 := InstanceAnalytics{
		Cpu:     10000,
		CpuPerc: 0.1,
	}
	instAnalytics2 := InstanceAnalytics{
		Cpu:     20000,
		CpuPerc: 0.2,
	}
	instAnalytics3 := InstanceAnalytics{
		Cpu:     30000,
		CpuPerc: 0.3,
	}
	InstAnalytics4 := InstanceAnalytics{
		Cpu:     0,
		CpuPerc: 0.0,
	}
	instances := map[string]InstanceAnalytics{
		"instance1": instAnalytics1,
		"instance2": instAnalytics2,
		"instance3": instAnalytics3,
		"instance4": InstAnalytics4,
	}

	system := SystemAnalytics{5000000}

	mockAnalytics := GruAnalytics{
		Service:  services,
		Instance: instances,
		System:   system,
	}

	gruAnalytics = mockAnalytics
}
