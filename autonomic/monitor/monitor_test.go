package monitor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveResource(t *testing.T) {
	mockStats := createMockStats()
	mockInstId := "instance3"

	removeResource(mockInstId, &mockStats)
	serviceStatsInst := mockStats.Service["service2"].Instances
	instancesStats := []string{}
	for k, _ := range mockStats.Instance {
		instancesStats = append(instancesStats, k)
	}

	assert.NotContains(t, serviceStatsInst, mockInstId, "Service stats should not contain 'instance3'")
	assert.NotContains(t, instancesStats, mockInstId, "Instance stats should not contain 'instance3'")
	assert.Contains(t, mockStats.Service["service2"].Events.Die, mockInstId, "Events Die should contain 'instance3'")

}

func TestFindServiceByInstanceId(t *testing.T) {
	mockStats := createMockStats()
	mockInstId := "instance3"

	mockService := findServiceByInstanceId(mockInstId, &mockStats)
	assert.Equal(t, "service2", mockService, "found service should be 'service2'")
}

func TestResetEventsStats(t *testing.T) {
	mockStats := createMockStats()
	srvName := "service1"

	resetEventsStats(srvName, &mockStats)
	assert.Equal(t, 0, len(mockStats.Service[srvName].Events.Die), "Events Die should be empty")
}

func TestCopyStats(t *testing.T) {
	mockStats := createMockStats()
	mockStats_cp := GruStats{
		Service:  make(map[string]ServiceStats),
		Instance: make(map[string]InstanceStats),
	}

	copyStats(&mockStats, &mockStats_cp)
	assert.Equal(t, mockStats.System.Cpu, mockStats_cp.System.Cpu, "Copy should be equal to the original")

	service := "service1"
	resetEventsStats(service, &mockStats)
	assert.Contains(t, mockStats_cp.Service[service].Events.Die,
		"instance0", "The copy should be not modified")
}

func createMockStats() GruStats {
	instances1 := []string{"instance1", "instance2"}
	events1 := EventStats{
		Die: []string{"instance0"},
	}
	service1 := ServiceStats{
		Instances: instances1,
		Events:    events1,
	}
	instances2 := []string{"instance3"}
	service2 := ServiceStats{Instances: instances2}
	services := map[string]ServiceStats{
		"service1": service1,
		"service2": service2,
	}

	instStat1 := InstanceStats{20000}
	instStat2 := InstanceStats{60000}
	instStat3 := InstanceStats{60000}
	instances := map[string]InstanceStats{
		"instance1": instStat1,
		"instance2": instStat2,
		"instance3": instStat3,
	}

	system := SystemStats{150000}

	mockStats := GruStats{
		Service:  services,
		Instance: instances,
		System:   system,
	}

	return mockStats
}

func TestFindIdIndex(t *testing.T) {
	instances := []string{
		"instance1",
		"instance2",
		"instance3",
		"instance4",
		"instance5",
	}

	index := findIdIndex("instance3", instances)
	assert.Equal(t, 2, index, "index of 'instance3' should be 2")

}
