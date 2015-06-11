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
	mockStats, names := createMockStats()
	for _, name := range names {
		updateInstances(name, &mockStats)
	}

	assert.NotContains(t, gruAnalytics.Service["service1"].Instances, "instance2", "Service 1 should contains instance2")
	inst := []string{}
	for k, _ := range gruAnalytics.Instance {
		inst = append(inst, k)
	}
	assert.NotContains(t, inst, "instance2", "Instances should not contains instance2")
	assert.Contains(t, gruAnalytics.Service["service2"].Instances, "instance3", "Service 2 should contains instance3")

	cleanAnalytics()
}

func TestUpdateAnalytics(t *testing.T) {
	mockStats, names := createMockStats()
	for _, name := range names {
		updateAnalytics(name, &mockStats)
	}
	statsCpu := mockStats.Instance["instance1"].Cpu
	analyticsCpu := gruAnalytics.Instance["instance1"].Cpu
	assert.Equal(t, statsCpu, analyticsCpu, "Instance1 stats and analytics should be equal")

	cleanAnalytics()
}

func TestComputeCpuAvg(t *testing.T) {
	mockStats, names := createMockStats()
	createMockAnalytics()
	for _, name := range names {
		computeCpuAvg(name, &mockStats)
	}

	assert.Equal(t, 0.25, gruAnalytics.Service["service1"].CpuAvg, "Service1 cpuAvg should be 25%")
	assert.Equal(t, 0.3, gruAnalytics.Service["service2"].CpuAvg, "Service2 cpuAvg should be 30%")
}

func createMockStats() (monitor.GruStats, []string) {
	instances1 := []string{"instance1"}
	events1 := monitor.EventStats{
		Die: []string{"instance2"},
	}
	service1 := monitor.ServiceStats{
		Instances: instances1,
		Events:    events1,
	}

	instances2 := []string{"instance3"}
	service2 := monitor.ServiceStats{Instances: instances2}
	services := map[string]monitor.ServiceStats{
		"service1": service1,
		"service2": service2,
	}

	instStat1 := monitor.InstanceStats{Cpu: 20000}
	instStat2 := monitor.InstanceStats{Cpu: 60000}
	instStat3 := monitor.InstanceStats{Cpu: 60000}
	instances := map[string]monitor.InstanceStats{
		"instance1": instStat1,
		"instance2": instStat2,
		"instance3": instStat3,
	}

	system := monitor.SystemStats{15000000}

	mockStats := monitor.GruStats{
		Service:  services,
		Instance: instances,
		System:   system,
	}

	return mockStats, []string{"service1", "service2"}
}

func cleanAnalytics() {
	gruAnalytics = GruAnalytics{}
	gruAnalytics.Service = make(map[string]ServiceAnalytics)
	gruAnalytics.Instance = make(map[string]InstanceAnalytics)
}

func createMockAnalytics() {
	instances1 := []string{"instance1", "instance2"}
	service1 := ServiceAnalytics{
		CpuAvg:    0.2,
		Instances: instances1}
	instances2 := []string{"instance3"}
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
	instances := map[string]InstanceAnalytics{
		"instance1": instAnalytics1,
		"instance2": instAnalytics2,
		"instance3": instAnalytics3,
	}

	system := SystemAnalytics{5000000}

	mockAnalytics := GruAnalytics{
		Service:  services,
		Instance: instances,
		System:   system,
	}

	gruAnalytics = mockAnalytics
}
