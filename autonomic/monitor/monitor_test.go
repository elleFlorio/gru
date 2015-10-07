package monitor

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/storage"
)

func init() {
	//Initialize storage
	datastore, _ := storage.New("internal")
	datastore.Initialize()
}

func TestUpdateRunningInstances(t *testing.T) {
	mockStats := CreateMockStats()
	history = CreateMockHistory()
	wsize := MaxNumberOfEntryInHistory()
	mockService := "service1"
	promoted := "instance1_3"

	updateRunningInstances(mockService, &mockStats, wsize)

	assert.Contains(t, mockStats.Service[mockService].Instances.Running, promoted,
		"(promoted) Service1 - instances - running should contain promoted instance")
}

func TestComputeInstanceCpuPerc(t *testing.T) {
	mockInstCpus := []float64{10000, 20000, 30000, 40000, 50000, 60000}
	mockSysCpus := []float64{1000000, 1100000, 1200000, 1300000, 1400000, 1500000}

	mockPerc := computeInstanceCpuPerc(mockInstCpus, mockSysCpus)

	assert.Equal(t, 0.1, mockPerc, "Computed percentage should be 1%")
}

func TestComputeServiceCpuPerc(t *testing.T) {
	mockStats := CreateMockStats()
	history = CreateMockHistory()
	srv1 := "service1"
	srv2 := "service2"

	computeServiceCpuPerc(srv1, &mockStats)
	computeServiceCpuPerc(srv2, &mockStats)
	cpuAvgS1 := mockStats.Service[srv1].Cpu.Avg
	cpuTotS1 := mockStats.Service[srv1].Cpu.Tot
	cpuAvgS2 := mockStats.Service[srv2].Cpu.Avg
	cpuTotS2 := mockStats.Service[srv2].Cpu.Tot

	assert.Equal(t, 0.35, cpuAvgS1, "Cpu average of service1 should be 0.35")
	assert.Equal(t, 0.7, cpuTotS1, "Cpu total of service1 should be 0.7")

	assert.Equal(t, 0.4, cpuAvgS2, "Cpu average of service2 should be 0.4")
	assert.Equal(t, 0.4, cpuTotS2, "Cpu total of service2 should be 0.4")

}

func TestMakeSnapshot(t *testing.T) {
	mockStats := CreateMockStats()
	history = CreateMockHistory()
	mockStats_cp := GruStats{
		Service:  make(map[string]ServiceStats),
		Instance: make(map[string]InstanceStats),
	}

	makeSnapshot(&mockStats, &mockStats_cp)
	assert.Equal(t,
		len(mockStats.Service["service1"].Instances.Running),
		len(mockStats_cp.Service["service1"].Instances.Running),
		"Copy should be equal to the original")
	assert.Equal(t,
		len(mockStats.System.Instances.Stopped),
		len(mockStats_cp.System.Instances.Stopped),
		"Copy should be equal to the original")

	service := "service1"
	resetEventsStats(service, &mockStats)
	assert.Contains(t, mockStats_cp.Service[service].Events.Stop,
		"instance1_0", "The copy should be not modified")
}

func TestResetEventsStats(t *testing.T) {
	mockStats := CreateMockStats()
	srvName := "service1"

	resetEventsStats(srvName, &mockStats)
	assert.Equal(t, 0, len(mockStats.Service[srvName].Events.Stop), "Events Stop should be empty")
}

func TestAddResource(t *testing.T) {
	mockStats := CreateMockStats()
	mockHist := CreateMockHistory()
	id2_s := "instance2_s"
	id2_p := "instance2_p"
	id2_r := "instance2_r"
	srvName := "service2"
	status2_s := "stopped"
	status2_p := "pending"
	status2_r := "running"

	// check add stopped
	addResource(id2_s, srvName, status2_s, &mockStats, &mockHist)
	assert.Contains(t, mockStats.Service[srvName].Instances.All, id2_s,
		"(new -> stopped) Service 2 - instances - all, should contain added instance")
	assert.Contains(t, mockStats.Service[srvName].Instances.Stopped, id2_s,
		"(new -> stopped) Service 2 - instances - stopped, should contain added instance")

	// check add pending
	addResource(id2_p, srvName, status2_p, &mockStats, &mockHist)
	assert.Contains(t, mockStats.Service[srvName].Instances.All, id2_p,
		"(new -> pending) Service 2 - instances - all, should contain added instance")
	assert.Contains(t, mockStats.Service[srvName].Instances.Pending, id2_p,
		"(new -> pending) Service 2 - instances - pending, should contain added instance")
	assert.Contains(t, mockStats.Service[srvName].Events.Start, id2_p,
		"(new -> pending) Service 2 - events - start, should contain added instance")

	// check add running
	addResource(id2_r, srvName, status2_r, &mockStats, &mockHist)
	assert.Contains(t, mockStats.Service[srvName].Instances.All, id2_r,
		"(new -> running) Service 2 - instances - all, should contain added instance")
	assert.Contains(t, mockStats.Service[srvName].Instances.Running, id2_r,
		"(new -> running) Service 2 - instances - running, should contain added instance")

	//check stopped -> pending
	addResource(id2_s, srvName, status2_p, &mockStats, &mockHist)
	assert.Contains(t, mockStats.Service[srvName].Instances.Pending, id2_s,
		"(stopped -> pending) Service 2 - instances - pending, should contain added instance")
	assert.Contains(t, mockStats.Service[srvName].Events.Start, id2_s,
		"(stopped -> pending) Service 2 - events - start, should contain added instance")
	assert.NotContains(t, mockStats.Service[srvName].Instances.Stopped, id2_s,
		"(stopped -> pending) Service 2 - instances - stopped, should not contain added instance")

	//check pending -> running
	addResource(id2_s, srvName, status2_r, &mockStats, &mockHist)
	assert.Contains(t, mockStats.Service[srvName].Instances.Running, id2_s,
		"(pending -> running) Service 2 - instances - running, should contain added instance")
	assert.NotContains(t, mockStats.Service[srvName].Instances.Pending, id2_s,
		"(pending -> running) Service 2 - instances - pending, should not contain added instance")
}

func TestRemoveResource(t *testing.T) {
	mockStats := CreateMockStats()
	mockHist := CreateMockHistory()
	mockInstId_r := "instance2_1"
	mockInstId_p := "instance1_3"

	// check running
	removeResource(mockInstId_r, &mockStats, &mockHist)
	serviceStatsInst := mockStats.Service["service2"].Instances.Running
	instancesStats := []string{}
	for k, _ := range mockStats.Instance {
		instancesStats = append(instancesStats, k)
	}

	assert.NotContains(t, serviceStatsInst, mockInstId_r,
		"(running) Service stats should not contain 'instance2_1'")
	assert.NotContains(t, instancesStats, mockInstId_r,
		"(running) Instance stats should not contain 'instance2_1'")
	assert.Contains(t, mockStats.Service["service2"].Events.Stop, mockInstId_r,
		"(running) Events Stop should contain 'instance2_1'")

	// check pending
	removeResource(mockInstId_p, &mockStats, &mockHist)
	serviceStatsInst = mockStats.Service["service2"].Instances.Pending
	instancesStats = []string{}
	for k, _ := range mockStats.Instance {
		instancesStats = append(instancesStats, k)
	}

	assert.NotContains(t, serviceStatsInst, mockInstId_p,
		"(pending) Service stats should not contain 'instance1_3'")
	assert.NotContains(t, instancesStats, mockInstId_r,
		"(running) Instance stats should not contain 'instance1_3'")
	assert.Contains(t, mockStats.Service["service2"].Events.Stop, mockInstId_r,
		"(running) Events Stop should contain 'instance1_3'")
}

func TestFindServiceByInstanceId(t *testing.T) {
	mockStats := CreateMockStats()
	mockInstId := "instance1_4"

	mockService := findServiceByInstanceId(mockInstId, &mockStats)
	assert.Equal(t, "service1", mockService, "found service should be 'service2'")
}

func TestFindIdIndex(t *testing.T) {
	instances := []string{
		"instance1_1",
		"instance1_2",
		"instance1_3",
		"instance1_4",
		"instance2_1",
	}

	index, _ := findIdIndex("instance1_3", instances)
	assert.Equal(t, 2, index, "index of 'instance3' should be 2")
}
