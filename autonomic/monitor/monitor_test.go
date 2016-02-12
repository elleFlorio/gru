package monitor

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/autonomic/monitor/logreader"
	cfg "github.com/elleFlorio/gru/configuration"
	res "github.com/elleFlorio/gru/resources"
	"github.com/elleFlorio/gru/service"
	"github.com/elleFlorio/gru/storage"
)

func init() {
	//Initialize storage
	storage.New("internal")
	metric.Manager().Start()
	res.CreateMockResources(1, "1G", 0, "0G")
	// n := cfg.Node{
	// 	Resources: cfg.NodeResources{
	// 		TotalCpus: 1,
	// 	},
	// }
	// cfg.SetNode(n)
	cfg.SetNode(cfg.Node{})
	resetMockServices()
}

func resetMockServices() {
	mockServices := service.CreateMockServices()
	cfg.SetServices(mockServices)
}

func TestUpdateRunningInstances(t *testing.T) {
	defer resetMockServices()

	mockStats := CreateMockStats()
	history = CreateMockHistory()
	wsize := MaxNumberOfEntryInHistory()
	mockService := "service1"
	promoted := "instance1_3"
	srv, _ := service.GetServiceByName(mockService)

	updateRunningInstances(mockService, &mockStats, wsize)

	assert.Contains(t, srv.Instances.Running, promoted)
}

func TestComputeInstanceCpuPerc(t *testing.T) {
	mockInstCpus := []float64{10000, 20000, 30000, 40000, 50000, 60000}
	mockSysCpus := []float64{1000000, 1100000, 1200000, 1300000, 1400000, 1500000}

	mockPerc := computeInstanceCpuPerc(mockInstCpus, mockSysCpus)

	assert.Equal(t, 0.1, mockPerc)
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

	assert.Equal(t, 0.35, cpuAvgS1)
	assert.Equal(t, 0.7, cpuTotS1)

	assert.Equal(t, 0.4, cpuAvgS2)
	assert.Equal(t, 0.4, cpuTotS2)

}

func TestUpdateSystemInstances(t *testing.T) {
	defer resetMockServices()

	mockStats := CreateMockStats()
	updateSystemInstances(&mockStats)

	srv1, _ := service.GetServiceByName("service1")
	srv2, _ := service.GetServiceByName("service2")
	tot_all := len(srv1.Instances.All) + len(srv2.Instances.All)
	tot_pen := len(srv1.Instances.Pending) + len(srv2.Instances.Pending)
	tot_run := len(srv1.Instances.Running) + len(srv2.Instances.Running)
	tot_stop := len(srv1.Instances.Stopped) + len(srv2.Instances.Stopped)
	tot_pause := len(srv1.Instances.Paused) + len(srv2.Instances.Paused)

	instances := cfg.GetNodeInstances()
	assert.Len(t, instances.All, tot_all)
	assert.Len(t, instances.Pending, tot_pen)
	assert.Len(t, instances.Running, tot_run)
	assert.Len(t, instances.Stopped, tot_stop)
	assert.Len(t, instances.Paused, tot_pause)
}

func TestMakeSnapshot(t *testing.T) {
	mockStats := CreateMockStats()
	history = CreateMockHistory()
	mockStats_cp := GruStats{
		Service:  make(map[string]ServiceStats),
		Instance: make(map[string]InstanceStats),
	}

	makeSnapshot(&mockStats, &mockStats_cp)
	service := "service1"
	resetEventsStats(service, &mockStats)
	assert.Contains(t, mockStats_cp.Service[service].Events.Stop, "instance1_0")
}

func TestResetEventsStats(t *testing.T) {
	mockStats := CreateMockStats()
	srvName := "service1"

	resetEventsStats(srvName, &mockStats)
	assert.Equal(t, 0, len(mockStats.Service[srvName].Events.Stop))
}

func TestAddResource(t *testing.T) {
	defer resetMockServices()

	mockStats := CreateMockStats()
	mockHist := CreateMockHistory()
	id2_s := "instance2_s"
	id2_p := "instance2_p"
	id2_r := "instance2_r"
	srvName := "service2"
	status2_s := "stopped"
	status2_p := "pending"
	status2_r := "running"

	srv, _ := service.GetServiceByName(srvName)
	// check add stopped
	addResource(id2_s, srvName, status2_s, &mockStats, &mockHist)
	assert.Contains(t, srv.Instances.All, id2_s,
		"(new -> stopped) Service 2 - instances - all, should contain added instance")
	assert.Contains(t, srv.Instances.Stopped, id2_s,
		"(new -> stopped) Service 2 - instances - stopped, should contain added instance")

	// check add pending
	addResource(id2_p, srvName, status2_p, &mockStats, &mockHist)
	assert.Contains(t, srv.Instances.All, id2_p,
		"(new -> pending) Service 2 - instances - all, should contain added instance")
	assert.Contains(t, srv.Instances.Pending, id2_p,
		"(new -> pending) Service 2 - instances - pending, should contain added instance")
	assert.Contains(t, mockStats.Service[srvName].Events.Start, id2_p,
		"(new -> pending) Service 2 - events - start, should contain added instance")

	// check add running
	addResource(id2_r, srvName, status2_r, &mockStats, &mockHist)
	assert.Contains(t, srv.Instances.All, id2_r,
		"(new -> running) Service 2 - instances - all, should contain added instance")
	assert.Contains(t, srv.Instances.Running, id2_r,
		"(new -> running) Service 2 - instances - running, should contain added instance")

	//check stopped -> pending
	addResource(id2_s, srvName, status2_p, &mockStats, &mockHist)
	assert.Contains(t, srv.Instances.Pending, id2_s,
		"(stopped -> pending) Service 2 - instances - pending, should contain added instance")
	assert.Contains(t, mockStats.Service[srvName].Events.Start, id2_s,
		"(stopped -> pending) Service 2 - events - start, should contain added instance")
	assert.NotContains(t, srv.Instances.Stopped, id2_s,
		"(stopped -> pending) Service 2 - instances - stopped, should not contain added instance")

	//check pending -> running
	addResource(id2_s, srvName, status2_r, &mockStats, &mockHist)
	assert.Contains(t, srv.Instances.Running, id2_s,
		"(pending -> running) Service 2 - instances - running, should contain added instance")
	assert.NotContains(t, srv.Instances.Pending, id2_s,
		"(pending -> running) Service 2 - instances - pending, should not contain added instance")
}

func TestRemoveResource(t *testing.T) {
	defer resetMockServices()

	mockStats := CreateMockStats()
	mockHist := CreateMockHistory()
	mockInstId_r := "instance2_1"
	mockInstId_p := "instance1_3"
	serviceName := "service2"
	srv, _ := service.GetServiceByName(serviceName)

	// check error
	removeResource("pippo", &mockStats, &mockHist)

	// check running
	removeResource(mockInstId_r, &mockStats, &mockHist)
	serviceStatsInst := srv.Instances.Running
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
	serviceStatsInst = srv.Instances.Pending
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

func TestConvertStatsToData(t *testing.T) {
	stats_ok := CreateMockStats()

	_, err := convertStatsToData(stats_ok)
	assert.NoError(t, err, "(ok) stats convertion should produce no error")
}

func TestConvertDataToStats(t *testing.T) {
	data_ok, err := convertStatsToData(CreateMockStats())
	data_bad := []byte{}

	_, err = convertDataToStats(data_ok)
	assert.NoError(t, err, "(ok) data convertion should produce no error")

	_, err = convertDataToStats(data_bad)
	assert.Error(t, err, "(bad) data convertion should produce an error")
}

func TestGetMonitorData(t *testing.T) {
	_, err := GetMonitorData()
	assert.Error(t, err)

	StoreMockStats()
	_, err = GetMonitorData()
	assert.NoError(t, err)
}

func TestRun(t *testing.T) {
	gruStats = CreateMockStats()
	history = CreateMockHistory()
	resetMockServices()

	assert.NotEmpty(t, Run())
}
