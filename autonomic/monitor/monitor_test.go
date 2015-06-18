package monitor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveResource(t *testing.T) {
	mockStats := CreateMockStats()
	mockInstId := "instance2_1"

	removeResource(mockInstId, &mockStats)
	serviceStatsInst := mockStats.Service["service2"].Instances.Running
	instancesStats := []string{}
	for k, _ := range mockStats.Instance {
		instancesStats = append(instancesStats, k)
	}

	assert.NotContains(t, serviceStatsInst, mockInstId, "Service stats should not contain 'instance3'")
	assert.NotContains(t, instancesStats, mockInstId, "Instance stats should not contain 'instance3'")
	assert.Contains(t, mockStats.Service["service2"].Events.Stop, mockInstId, "Events Stop should contain 'instance3'")

}

func TestFindServiceByInstanceId(t *testing.T) {
	mockStats := CreateMockStats()
	mockInstId := "instance1_4"

	mockService := findServiceByInstanceId(mockInstId, &mockStats)
	assert.Equal(t, "service1", mockService, "found service should be 'service2'")
}

func TestResetEventsStats(t *testing.T) {
	mockStats := CreateMockStats()
	srvName := "service1"

	resetEventsStats(srvName, &mockStats)
	assert.Equal(t, 0, len(mockStats.Service[srvName].Events.Stop), "Events Stop should be empty")
}

func TestCopyStats(t *testing.T) {
	mockStats := CreateMockStats()
	mockStats_cp := GruStats{
		Service:  make(map[string]ServiceStats),
		Instance: make(map[string]InstanceStats),
	}

	copyStats(&mockStats, &mockStats_cp)
	assert.Equal(t, mockStats.System.Cpu, mockStats_cp.System.Cpu, "Copy should be equal to the original")

	service := "service1"
	resetEventsStats(service, &mockStats)
	assert.Contains(t, mockStats_cp.Service[service].Events.Stop,
		"instance1_0", "The copy should be not modified")
}

func TestFindIdIndex(t *testing.T) {
	instances := []string{
		"instance1_1",
		"instance1_2",
		"instance1_3",
		"instance1_4",
		"instance2_1",
	}

	index := findIdIndex("instance1_3", instances)
	assert.Equal(t, 2, index, "index of 'instance3' should be 2")

}
