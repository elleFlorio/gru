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
	defer cleanAnalytics()
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
}

func TestUpdateAnalytics(t *testing.T) {
	defer cleanAnalytics()
	mockStats := createMockStats()
	for name, _ := range mockStats.Service {
		updateAnalytics(name, &mockStats)
	}
	statsCpu := mockStats.Instance["instance1"].Cpu
	analyticsCpu := gruAnalytics.Instance["instance1"].Cpu
	assert.Equal(t, statsCpu, analyticsCpu, "Instance1 stats and analytics should be equal")
}

func TestComputeCpuAvg(t *testing.T) {
	defer cleanAnalytics()
	var err error = nil
	mockStats := createMockStats()
	createMockAnalytics()

	err = computeCpuAvg("service1", &mockStats)
	assert.Equal(t, 0.25, gruAnalytics.Service["service1"].CpuAvg, "Service1 cpuAvg should be 25%")
	err = computeCpuAvg("service2", &mockStats)
	assert.Equal(t, 0.3, gruAnalytics.Service["service2"].CpuAvg, "Service2 cpuAvg should be 30%")
	err = computeCpuAvg("service3", &mockStats)
	assert.Error(t, err, "Service 3 cpu avg should be not valid")
	assert.Equal(t, 0.0, gruAnalytics.Service["service3"].CpuAvg, "Service3 cpuAvg should be 0%")
}

func cleanAnalytics() {
	gruAnalytics = GruAnalytics{}
	gruAnalytics.Service = make(map[string]ServiceAnalytics)
	gruAnalytics.Instance = make(map[string]InstanceAnalytics)
}

func createMockAnalytics() {
	all1 := []string{"instance0", "instance1", "instance2", "instance4"}
	pending1 := []string{"instance4"}
	running1 := []string{"instance1", "instance2"}
	instances1 := InstanceStatus{
		All:     all1,
		Pending: pending1,
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

	all3 := []string{"instance5"}
	pending3 := []string{"instance5"}
	instances3 := InstanceStatus{
		All:     all3,
		Pending: pending3,
	}
	service3 := ServiceAnalytics{
		CpuAvg:    0.0,
		Instances: instances3,
	}

	services := map[string]ServiceAnalytics{
		"service1": service1,
		"service2": service2,
		"service3": service3,
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
	InstAnalytics5 := InstanceAnalytics{
		Cpu:     0,
		CpuPerc: 0.0,
	}
	instances := map[string]InstanceAnalytics{
		"instance1": instAnalytics1,
		"instance2": instAnalytics2,
		"instance3": instAnalytics3,
		"instance4": InstAnalytics4,
		"instance5": InstAnalytics5,
	}

	system := SystemAnalytics{5000000}

	mockAnalytics := GruAnalytics{
		Service:  services,
		Instance: instances,
		System:   system,
	}

	gruAnalytics = mockAnalytics
}
