package monitor

import (
	"testing"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/samalba/dockerclient"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/autonomic/monitor/logreader"
	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/data"
	"github.com/elleFlorio/gru/discovery"
	res "github.com/elleFlorio/gru/resources"
	"github.com/elleFlorio/gru/service"
	"github.com/elleFlorio/gru/storage"
)

func init() {
	log.SetLevel(log.ErrorLevel)
	storage.New("internal")
	discovery.New("noservice", "")
	cfg.GetAgentDiscovery().TTL = 5
	metric.Manager().Start()
	res.CreateMockResources(1, "1G", 0, "0G")
	cfg.SetNode(cfg.Node{})
	resetMockServices()
}

func TestUpdateRunningInstances(t *testing.T) {
	defer resetMockServices()

	mockStats := data.CreateMockStats()
	history = data.CreateMockHistory()
	wsize := data.MaxNumberOfEntryInHistory()
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

	mockInstCpus = []float64{10000, 10000, 10000, 10000, 10000, 10000}

	mockPerc = computeInstanceCpuPerc(mockInstCpus, mockSysCpus)
	assert.Equal(t, 0.0, mockPerc)
}

func TestComputeServiceCpuPerc(t *testing.T) {
	mockStats := data.CreateMockStats()
	history = data.CreateMockHistory()
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

	mockStats := data.CreateMockStats()
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
	mockStats := data.CreateMockStats()
	history = data.CreateMockHistory()
	mockStats_cp := data.GruStats{
		Service:  make(map[string]data.ServiceStats),
		Instance: make(map[string]data.InstanceStats),
	}

	makeSnapshot(&mockStats, &mockStats_cp)
	service := "service1"
	resetEventsStats(service, &mockStats)
	assert.Contains(t, mockStats_cp.Service[service].Events.Stop, "instance1_0")
}

func TestResetEventsStats(t *testing.T) {
	mockStats := data.CreateMockStats()
	srvName := "service1"

	resetEventsStats(srvName, &mockStats)
	assert.Equal(t, 0, len(mockStats.Service[srvName].Events.Stop))
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

func TestAddInstance(t *testing.T) {
	defer resetMockServices()

	mockStats := data.CreateMockStats()
	mockHist := data.CreateMockHistory()
	id2_s := "instance2_s"
	id2_p := "instance2_p"
	id2_r := "instance2_r"
	srvName := "service2"
	status2_s := "stopped"
	status2_p := "pending"
	status2_r := "running"

	srv, _ := service.GetServiceByName(srvName)
	// check add stopped
	addInstance(id2_s, srvName, status2_s, &mockStats, &mockHist)
	assert.Contains(t, srv.Instances.All, id2_s,
		"(new -> stopped) Service 2 - instances - all, should contain added instance")
	assert.Contains(t, srv.Instances.Stopped, id2_s,
		"(new -> stopped) Service 2 - instances - stopped, should contain added instance")

	// check add pending
	addInstance(id2_p, srvName, status2_p, &mockStats, &mockHist)
	assert.Contains(t, srv.Instances.All, id2_p,
		"(new -> pending) Service 2 - instances - all, should contain added instance")
	assert.Contains(t, srv.Instances.Pending, id2_p,
		"(new -> pending) Service 2 - instances - pending, should contain added instance")
	assert.Contains(t, mockStats.Service[srvName].Events.Start, id2_p,
		"(new -> pending) Service 2 - events - start, should contain added instance")

	// check add running
	addInstance(id2_r, srvName, status2_r, &mockStats, &mockHist)
	assert.Contains(t, srv.Instances.All, id2_r,
		"(new -> running) Service 2 - instances - all, should contain added instance")
	assert.Contains(t, srv.Instances.Running, id2_r,
		"(new -> running) Service 2 - instances - running, should contain added instance")

	//check stopped -> pending
	addInstance(id2_s, srvName, status2_p, &mockStats, &mockHist)
	assert.Contains(t, srv.Instances.Pending, id2_s,
		"(stopped -> pending) Service 2 - instances - pending, should contain added instance")
	assert.Contains(t, mockStats.Service[srvName].Events.Start, id2_s,
		"(stopped -> pending) Service 2 - events - start, should contain added instance")
	assert.NotContains(t, srv.Instances.Stopped, id2_s,
		"(stopped -> pending) Service 2 - instances - stopped, should not contain added instance")

	//check pending -> running
	addInstance(id2_s, srvName, status2_r, &mockStats, &mockHist)
	assert.Contains(t, srv.Instances.Running, id2_s,
		"(pending -> running) Service 2 - instances - running, should contain added instance")
	assert.NotContains(t, srv.Instances.Pending, id2_s,
		"(pending -> running) Service 2 - instances - pending, should not contain added instance")
}

func TeststopInstance(t *testing.T) {
	defer resetMockServices()

	mockStats := data.CreateMockStats()
	mockHist := data.CreateMockHistory()
	mockInstId_r := "instance2_1"
	mockInstId_p := "instance1_3"
	serviceName := "service2"
	srv, _ := service.GetServiceByName(serviceName)

	// check error
	stopInstance("pippo", &mockStats, &mockHist)

	// check running
	stopInstance(mockInstId_r, &mockStats, &mockHist)
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
	stopInstance(mockInstId_p, &mockStats, &mockHist)
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

func TestRemoveInstance(t *testing.T) {
	defer resetMockServices()
	defer log.SetLevel(log.ErrorLevel)

	mockStats := data.CreateMockStats()
	mockHist := data.CreateMockHistory()
	service1 := "service1"
	mockInstId_s := "instance1_0"
	service2 := "service2"
	mockInstId_r := "instance2_1"
	mockInstId_wrong := "pippo"

	removeInstance(mockInstId_s, &mockStats, &mockHist)
	srv1, _ := service.GetServiceByName(service1)
	assert.NotContains(t, srv1.Instances.Running, mockInstId_r)
	assert.NotContains(t, srv1.Instances.Pending, mockInstId_r)
	assert.NotContains(t, srv1.Instances.Stopped, mockInstId_r)

	removeInstance(mockInstId_r, &mockStats, &mockHist)
	srv2, _ := service.GetServiceByName(service2)
	assert.NotContains(t, srv2.Instances.Running, mockInstId_r)
	assert.NotContains(t, srv2.Instances.Pending, mockInstId_r)
	assert.NotContains(t, srv2.Instances.Stopped, mockInstId_r)

	// Check the log for this test
	log.SetLevel(log.DebugLevel)
	removeInstance(mockInstId_wrong, &mockStats, &mockHist)

}

// func TestConvertStatsToData(t *testing.T) {
// 	stats_ok := data.CreateMockStats()

// 	_, err := convertStatsToData(stats_ok)
// 	assert.NoError(t, err, "(ok) stats convertion should produce no error")
// }

// func TestConvertDataToStats(t *testing.T) {
// 	data_ok, err := convertStatsToData(data.CreateMockStats())
// 	data_bad := []byte{}

// 	_, err = convertDataToStats(data_ok)
// 	assert.NoError(t, err, "(ok) data convertion should produce no error")

// 	_, err = convertDataToStats(data_bad)
// 	assert.Error(t, err, "(bad) data convertion should produce an error")
// }

// func TestGetMonitorData(t *testing.T) {
// 	_, err := GetMonitorData()
// 	assert.Error(t, err)

// 	StoreMockStats()
// 	_, err = GetMonitorData()
// 	assert.NoError(t, err)
// }

func TestRun(t *testing.T) {
	gruStats = data.CreateMockStats()
	history = data.CreateMockHistory()
	resetMockServices()

	assert.NotEmpty(t, Run())
}

func resetMockServices() {
	mockServices := service.CreateMockServices()
	cfg.SetServices(mockServices)
}

func TestCreatePortBindings(t *testing.T) {
	dockerBindings := createDockerBindings()
	portBindings := createPortBindings(dockerBindings)

	assert.Len(t, portBindings, 2)
	assert.Equal(t, "50100", portBindings["50100"][0])
	assert.Equal(t, "50200", portBindings["50200"][0])

}

func createDockerBindings() map[string][]dockerclient.PortBinding {
	port1 := "50100/tcp"
	host1 := dockerclient.PortBinding{
		HostIp:   "0.0.0.0",
		HostPort: "50100",
	}
	bindings1 := []dockerclient.PortBinding{host1}
	port2 := "50200/tcp"
	host2 := dockerclient.PortBinding{
		HostIp:   "0.0.0.0",
		HostPort: "50200",
	}
	bindings2 := []dockerclient.PortBinding{host2}

	dockerBindings := map[string][]dockerclient.PortBinding{
		port1: bindings1,
		port2: bindings2,
	}

	return dockerBindings
}
