package analyzer

import (
	"math/rand"
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/autonomic/monitor"
	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/node"
	"github.com/elleFlorio/gru/service"
	"github.com/elleFlorio/gru/storage"
	"github.com/elleFlorio/gru/utils"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const c_EPSILON = 0.09

func init() {
	gruAnalytics = GruAnalytics{
		Service: make(map[string]ServiceAnalytics),
	}

	//Initialize storage
	storage.New("internal")
}

//TODO this is not so straithforward
func TestComputeServiceResources(t *testing.T) {
	cpu8 := "0,1,2,3,4,5,6,7"
	cpu6 := "0,1,2,3,4,5"
	cpu4 := "0,1,2,3"
	cpu2 := "0,1"
	cpu1 := "0"

	n_full := createNode(6, "8G", 6, "8G")
	n_half_full := createNode(6, "8G", 4, "4G")
	n_half_empty := createNode(6, "8G", 2, "2G")
	n_empty := createNode(6, "8G", 0, "0G")

	name := "test"
	s_over := createService(name, cpu8, "16G")
	s_bigger := createService(name, cpu6, "8G")
	s_big := createService(name, cpu4, "4G")
	s_medium := createService(name, cpu2, "4G")
	s_low := createService(name, cpu2, "2G")
	s_lower := createService(name, cpu1, "1G")
	s_error := createService(name, cpu1, "error")

	//error
	node.UpdateNodeConfig(n_empty)
	service.UpdateServices([]service.Service{s_error})
	assert.Equal(t, 0.0, resourcesAvailable(name))

	// Node is full, all labels should be red
	node.UpdateNodeConfig(n_full)
	service.UpdateServices([]service.Service{s_over})
	assert.Equal(t, 0.0, resourcesAvailable(name))
	service.UpdateServices([]service.Service{s_bigger})
	assert.Equal(t, 0.0, resourcesAvailable(name))
	service.UpdateServices([]service.Service{s_big})
	assert.Equal(t, 0.0, resourcesAvailable(name))
	service.UpdateServices([]service.Service{s_medium})
	assert.Equal(t, 0.0, resourcesAvailable(name))
	service.UpdateServices([]service.Service{s_low})
	assert.Equal(t, 0.0, resourcesAvailable(name))
	service.UpdateServices([]service.Service{s_lower})
	assert.Equal(t, 0.0, resourcesAvailable(name))

	node.UpdateNodeConfig(n_half_full)
	service.UpdateServices([]service.Service{s_over})
	assert.Equal(t, 0.0, resourcesAvailable(name))
	service.UpdateServices([]service.Service{s_bigger})
	assert.Equal(t, 0.0, resourcesAvailable(name))
	service.UpdateServices([]service.Service{s_big})
	assert.Equal(t, 0.0, resourcesAvailable(name))
	service.UpdateServices([]service.Service{s_medium})
	assert.Equal(t, 1.0, resourcesAvailable(name))
	service.UpdateServices([]service.Service{s_low})
	assert.Equal(t, 1.0, resourcesAvailable(name))
	service.UpdateServices([]service.Service{s_lower})
	assert.Equal(t, 1.0, resourcesAvailable(name))

	node.UpdateNodeConfig(n_half_empty)
	service.UpdateServices([]service.Service{s_over})
	assert.Equal(t, 0.0, resourcesAvailable(name))
	service.UpdateServices([]service.Service{s_bigger})
	assert.Equal(t, 0.0, resourcesAvailable(name))
	service.UpdateServices([]service.Service{s_big})
	assert.Equal(t, 1.0, resourcesAvailable(name))
	service.UpdateServices([]service.Service{s_medium})
	assert.Equal(t, 1.0, resourcesAvailable(name))
	service.UpdateServices([]service.Service{s_low})
	assert.Equal(t, 1.0, resourcesAvailable(name))
	service.UpdateServices([]service.Service{s_lower})
	assert.Equal(t, 1.0, resourcesAvailable(name))

	node.UpdateNodeConfig(n_empty)
	service.UpdateServices([]service.Service{s_over})
	assert.Equal(t, 0.0, resourcesAvailable(name))
	service.UpdateServices([]service.Service{s_bigger})
	assert.Equal(t, 1.0, resourcesAvailable(name))
	service.UpdateServices([]service.Service{s_big})
	assert.Equal(t, 1.0, resourcesAvailable(name))
	service.UpdateServices([]service.Service{s_medium})
	assert.Equal(t, 1.0, resourcesAvailable(name))
	service.UpdateServices([]service.Service{s_low})
	assert.Equal(t, 1.0, resourcesAvailable(name))
	service.UpdateServices([]service.Service{s_lower})
	assert.Equal(t, 1.0, resourcesAvailable(name))
}

// CpuSet or CpuShares?
func createService(name string, cpu string, mem string) service.Service {
	srvConfig := service.Config{
		CpusetCpus: cpu,
		Memory:     mem,
	}

	srv := service.Service{
		Name:          name,
		Configuration: srvConfig,
	}

	return srv
}

func createNode(totCpu int64, totMem string, usedCpu int64, usedMem string) node.Node {
	totMemB, _ := utils.RAMInBytes(totMem)
	usedMemB, _ := utils.RAMInBytes(usedMem)

	resources := node.Resources{totMemB, totCpu, usedMemB, usedCpu}
	nd := node.Node{Resources: resources}

	return nd
}

func TestAnalyzeServices(t *testing.T) {
	analytics := GruAnalytics{
		Service: make(map[string]ServiceAnalytics),
	}

	services := []service.Service{
		createService("service1", "0,1,2,3", "4G"),
		createService("service2", "0", "2G"),
	}
	service.UpdateServices(services)

	nd := createNode(6, "8G", 0, "0G")
	node.UpdateNodeConfig(nd)

	stats := monitor.CreateMockStats()
	analyzeServices(&analytics, stats)

	assert.Len(t,
		analytics.Service["service1"].Instances.Running,
		len(stats.Service["service1"].Instances.Running),
	)
}

func TestComputeSystemResources(t *testing.T) {
	n_full := createNode(6, "8G", 6, "8G")
	n_half_full := createNode(6, "8G", 4, "6G")
	n_half := createNode(6, "8G", 3, "4G")
	n_half_empty := createNode(6, "8G", 2, "2G")
	n_empty := createNode(6, "8G", 0, "0G")

	node.UpdateNodeConfig(n_full)
	assert.Equal(t, 0.0, systemResourcesAvailable())
	node.UpdateNodeConfig(n_half_full)
	assert.InEpsilon(t, 0.3, systemResourcesAvailable(), c_EPSILON)
	node.UpdateNodeConfig(n_half)
	assert.InEpsilon(t, 0.5, systemResourcesAvailable(), c_EPSILON)
	node.UpdateNodeConfig(n_half_empty)
	assert.InEpsilon(t, 0.7, systemResourcesAvailable(), c_EPSILON)
	node.UpdateNodeConfig(n_empty)
	assert.InEpsilon(t, 1.0, systemResourcesAvailable(), c_EPSILON)
}

func TestAnalyzeSystem(t *testing.T) {
	stats := monitor.CreateMockStats()
	analyzeSystem(&gruAnalytics, stats)

	assert.InEpsilon(t, 1.0, gruAnalytics.System.Health, c_EPSILON)
}

func TestComputeNodeHealth(t *testing.T) {
	servicesH := []float64{
		1.0,
		1.0,
		0.6,
		0.8,
		0.4,
	}

	systemH := 0.8

	analytics := createHealth(servicesH, systemH)
	computeNodeHealth(&analytics)

	assert.Equal(t, 0.78, analytics.Health)
}

func createHealth(servicesH []float64, systemH float64) GruAnalytics {
	name := 'a'
	analytics := GruAnalytics{
		Service: make(map[string]ServiceAnalytics),
	}

	for _, h := range servicesH {
		srvA := ServiceAnalytics{Health: h}
		analytics.Service[string(name)] = srvA
		name += 1
	}

	analytics.System.Health = systemH

	return analytics

}

func TestGetPeerAnalytics(t *testing.T) {
	defer storage.DeleteAllData(enum.ANALYTICS)
	peer_1 := CreateMockAnalytics()
	peer_2 := CreateMockAnalytics()
	peer_3 := CreateMockAnalytics()

	data_1, _ := convertAnalyticsToData(peer_1)
	data_2, _ := convertAnalyticsToData(peer_2)
	data_3, _ := convertAnalyticsToData(peer_3)

	storage.StoreData("peer1", data_1, enum.ANALYTICS)
	storage.StoreData("peer2", data_2, enum.ANALYTICS)
	storage.StoreData("peer3", data_3, enum.ANALYTICS)

	peers := getPeersAnalytics()

	assert.Len(t, peers, 3)
}

func TestComputeServicesAvg(t *testing.T) {
	service.UpdateServices(createServicesWithNames([]string{"s1", "s2", "s3", "s4"}))
	analytics := createLocal()
	peers := createPeers()

	computeServicesAvg(peers, &analytics)
	//SERVICE 1
	assert.InEpsilon(t, 0.36, analytics.Service["s1"].Load, c_EPSILON)
	assert.InEpsilon(t, 0.28, analytics.Service["s1"].Resources.Cpu, c_EPSILON)
	assert.InEpsilon(t, 0.28, analytics.Service["s1"].Resources.Memory, c_EPSILON)
	assert.InEpsilon(t, 0.28, analytics.Service["s1"].Health, c_EPSILON)
	assert.Len(t, analytics.Service["s1"].Instances.All, 12)
	assert.Len(t, analytics.Service["s1"].Instances.Running, 5)
	assert.Len(t, analytics.Service["s1"].Instances.Pending, 2)
	assert.Len(t, analytics.Service["s1"].Instances.Stopped, 4)
	assert.Len(t, analytics.Service["s1"].Instances.Paused, 1)
	//SERVICE 2
	assert.Equal(t, 0.4, analytics.Service["s2"].Load)
	assert.Equal(t, 0.6, analytics.Service["s2"].Resources.Cpu)
	assert.Equal(t, 0.2, analytics.Service["s2"].Resources.Memory)
	assert.Equal(t, 0.6, analytics.Service["s2"].Health)
	assert.Len(t, analytics.Service["s2"].Instances.All, 3)
	assert.Len(t, analytics.Service["s2"].Instances.Running, 1)
	assert.Len(t, analytics.Service["s2"].Instances.Pending, 1)
	assert.Len(t, analytics.Service["s2"].Instances.Stopped, 0)
	assert.Len(t, analytics.Service["s2"].Instances.Paused, 1)
	//SERVICE 3
	assert.Equal(t, 0.6, analytics.Service["s3"].Load)
	assert.Equal(t, 0.6, analytics.Service["s3"].Resources.Cpu)
	assert.Equal(t, 0.6, analytics.Service["s3"].Resources.Memory)
	assert.Equal(t, 0.6, analytics.Service["s3"].Health)
	assert.Len(t, analytics.Service["s3"].Instances.All, 2)
	assert.Len(t, analytics.Service["s3"].Instances.Running, 1)
	assert.Len(t, analytics.Service["s3"].Instances.Pending, 1)
	assert.Len(t, analytics.Service["s3"].Instances.Stopped, 0)
	assert.Len(t, analytics.Service["s3"].Instances.Paused, 0)
	//SERVICE 4
	assert.Equal(t, 1.0, analytics.Service["s4"].Load)
	assert.Equal(t, 1.0, analytics.Service["s4"].Resources.Cpu)
	assert.Equal(t, 1.0, analytics.Service["s4"].Resources.Memory)
	assert.Equal(t, 1.0, analytics.Service["s4"].Health)
	assert.Len(t, analytics.Service["s4"].Instances.All, 2)
	assert.Len(t, analytics.Service["s4"].Instances.Running, 2)
	assert.Len(t, analytics.Service["s4"].Instances.Pending, 0)
	assert.Len(t, analytics.Service["s4"].Instances.Stopped, 0)
	assert.Len(t, analytics.Service["s4"].Instances.Paused, 0)

}

func createServicesWithNames(names []string) []service.Service {
	srvcs := []service.Service{}
	for _, name := range names {
		srv := service.Service{
			Name: name,
		}
		srvcs = append(srvcs, srv)
	}

	return srvcs
}

func createLocal() GruAnalytics {
	local := GruAnalytics{
		Service: make(map[string]ServiceAnalytics),
	}

	s1_is := createInstaceStatus(1, 0, 1, 1)
	s3_is := createInstaceStatus(1, 1, 0, 0)
	s4_is := createInstaceStatus(2, 0, 0, 0)
	s1_sa := createServiceAnalytics(0.2, 0.2, 0.2, 0.2, s1_is)
	s3_sa := createServiceAnalytics(0.6, 0.6, 0.6, 0.6, s3_is)
	s4_sa := createServiceAnalytics(1.0, 1.0, 1.0, 1.0, s4_is)
	local.Service["s1"] = s1_sa
	local.Service["s3"] = s3_sa
	local.Service["s4"] = s4_sa

	local.System = createSystemAnalytics(
		[]string{"s1", "s3", "s4"},
		0.8,
		0.6,
		0.6,
		s1_is, s3_is, s4_is)

	return local

}

func createPeers() []GruAnalytics {
	p1 := GruAnalytics{
		Service: make(map[string]ServiceAnalytics),
	}
	p2 := GruAnalytics{
		Service: make(map[string]ServiceAnalytics),
	}

	s1_is := createInstaceStatus(1, 2, 2, 0)
	s1b_is := createInstaceStatus(3, 0, 1, 0)
	s2_is := createInstaceStatus(1, 1, 0, 1)
	s1_sa := createServiceAnalytics(0.6, 0.4, 0.4, 0.4, s1_is)
	s1b_sa := createServiceAnalytics(1.0, 0.8, 0.8, 0.8, s1b_is)
	s2_sa := createServiceAnalytics(0.4, 0.6, 0.2, 0.6, s2_is)
	p1.Service["s1"] = s1_sa
	p2.Service["s1"] = s1b_sa
	p2.Service["s2"] = s2_sa

	p1.System = createSystemAnalytics(
		[]string{"s1"},
		0.6,
		0.6,
		0.6,
		s1_is)
	p2.System = createSystemAnalytics(
		[]string{"s1", "s2"},
		0.8,
		1.0,
		1.0,
		s1b_is, s2_is)

	p1.Health = 0.4
	p2.Health = 1.0

	return []GruAnalytics{p1, p2}
}

func createInstaceStatus(nInstRun int, nInstPend int, nInstStop int, nInstPaus int) service.InstanceStatus {
	nInstAll := nInstRun + nInstPend + nInstStop + nInstPaus

	inst_all := createRandomInstanceNames(nInstAll, 5)
	inst_run := createRandomInstanceNames(nInstRun, 5)
	inst_pen := createRandomInstanceNames(nInstPend, 5)
	inst_stp := createRandomInstanceNames(nInstStop, 5)
	inst_pau := createRandomInstanceNames(nInstPaus, 5)

	instStat := service.InstanceStatus{
		inst_all,
		inst_run,
		inst_pen,
		inst_stp,
		inst_pau,
	}

	return instStat
}

func createRandomInstanceNames(n int, length int) []string {
	names := []string{}
	for i := 0; i < n; i++ {
		names = append(names, randStringBytes(length))
	}
	return names
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func createServiceAnalytics(load float64, cpu float64, mem float64,
	health float64, instStatus service.InstanceStatus) ServiceAnalytics {

	srvA := ServiceAnalytics{
		Load: load,
		Resources: ResourcesAnalytics{
			Cpu:    cpu,
			Memory: mem,
		},
		Instances: instStatus,
		Health:    health,
	}

	return srvA
}

func createSystemAnalytics(services []string, cpu float64,
	mem float64, health float64, instStatus ...service.InstanceStatus) SystemAnalytics {

	sysInstSt := service.InstanceStatus{}
	for _, st := range instStatus {
		sysInstSt.All = append(sysInstSt.All, st.All...)
		sysInstSt.Pending = append(sysInstSt.Pending, st.Pending...)
		sysInstSt.Running = append(sysInstSt.Running, st.Running...)
		sysInstSt.Stopped = append(sysInstSt.Stopped, st.Stopped...)
		sysInstSt.Paused = append(sysInstSt.Paused, st.Paused...)
	}

	sysA := SystemAnalytics{
		Services: services,
		Resources: ResourcesAnalytics{
			Cpu:    cpu,
			Memory: mem,
		},
		Instances: sysInstSt,
		Health:    health,
	}

	return sysA

}

func TestCheckAndAppend(t *testing.T) {
	slice1 := []string{"a", "b", "c"}
	slice2 := []string{"b", "c", "d"}

	assert.Len(t, checkAndAppend(slice1, slice2), 4)
}

func TestComputeClusterAvg(t *testing.T) {
	analytics := createLocal()
	peers := createPeers()

	computeClusterAvg(peers, &analytics)
	assert.Len(t, analytics.Cluster.Services, 4)
	assert.InEpsilon(t, 0.73, analytics.Cluster.ResourcesAnalytics.Cpu, c_EPSILON)
	assert.InEpsilon(t, 0.73, analytics.Cluster.ResourcesAnalytics.Memory, c_EPSILON)
	assert.InEpsilon(t, 0.73, analytics.Cluster.Health, c_EPSILON)
}

func TestAnalyzeCluster(t *testing.T) {
	analytics := createLocal()
	defer storage.DeleteAllData(enum.ANALYTICS)
	peer_1 := CreateMockAnalytics()
	peer_2 := CreateMockAnalytics()
	peer_3 := CreateMockAnalytics()

	data_1, _ := convertAnalyticsToData(peer_1)
	data_2, _ := convertAnalyticsToData(peer_2)
	data_3, _ := convertAnalyticsToData(peer_3)

	storage.StoreData("peer1", data_1, enum.ANALYTICS)
	storage.StoreData("peer2", data_2, enum.ANALYTICS)
	storage.StoreData("peer3", data_3, enum.ANALYTICS)

	analyzeCluster(&analytics)
}

func TestSaveAnalytics(t *testing.T) {
	defer storage.DeleteAllData(enum.ANALYTICS)
	var err error
	analytics := CreateMockAnalytics()
	err = saveAnalytics(analytics)
	assert.NoError(t, err)
}

func TestGetAnalyzerData(t *testing.T) {
	defer storage.DeleteAllData(enum.ANALYTICS)
	var err error

	_, err = GetAnalyzerData()
	assert.Error(t, err)

	StoreMockAnalytics()
	_, err = GetAnalyzerData()
	assert.NoError(t, err)
}
