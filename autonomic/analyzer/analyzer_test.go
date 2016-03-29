package analyzer

import (
	"math/rand"
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/data"
	"github.com/elleFlorio/gru/enum"
	res "github.com/elleFlorio/gru/resources"
	"github.com/elleFlorio/gru/service"
	"github.com/elleFlorio/gru/storage"
	"github.com/elleFlorio/gru/utils"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const c_EPSILON = 0.09

func init() {
	gruAnalytics = data.GruAnalytics{
		Service: make(map[string]data.ServiceAnalytics),
	}

	//Initialize storage
	storage.New("internal")
}

func TestAnalyzeServices(t *testing.T) {
	analytics := data.GruAnalytics{
		Service: make(map[string]data.ServiceAnalytics),
	}

	services := []cfg.Service{
		createService("service1", "0,1,2,3", "4G", []string{"a", "b", "c"}),
		createService("service2", "0", "2G", []string{"d"}),
	}
	cfg.SetServices(services)
	setResources(6, "8G", 0, "0G")
	stats := data.CreateMockStats()

	analyzeServices(&analytics, stats)

	srv, _ := service.GetServiceByName("service1")
	assert.Len(t,
		analytics.Service["service1"].Instances.Running,
		len(srv.Instances.Running),
	)
}

func createService(name string, cpu string, mem string, running []string) cfg.Service {
	srvConfig := cfg.ServiceDocker{
		CpusetCpus: cpu,
		Memory:     mem,
	}

	srvStatus := cfg.ServiceStatus{Running: running}

	srv := cfg.Service{
		Name:      name,
		Docker:    srvConfig,
		Instances: srvStatus,
	}

	return srv
}

func setResources(totCpu int64, totMem string, usedCpu int64, usedMem string) {
	totMemB, _ := utils.RAMInBytes(totMem)
	usedMemB, _ := utils.RAMInBytes(usedMem)
	resources := res.GetResources()
	resources.Memory.Total = totMemB
	resources.Memory.Used = usedMemB
	resources.CPU.Total = totCpu
	resources.CPU.Used = usedCpu
}

func TestAnalyzeSystem(t *testing.T) {
	stats := data.CreateMockStats()
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

func createHealth(servicesH []float64, systemH float64) data.GruAnalytics {
	name := 'a'
	analytics := data.GruAnalytics{
		Service: make(map[string]data.ServiceAnalytics),
	}

	for _, h := range servicesH {
		srvA := data.ServiceAnalytics{Health: h}
		analytics.Service[string(name)] = srvA
		name += 1
	}

	analytics.System.Health = systemH

	return analytics

}

func createServicesWithNames(names []string) []cfg.Service {
	srvcs := []cfg.Service{}
	for _, name := range names {
		srv := cfg.Service{
			Name: name,
		}
		srvcs = append(srvcs, srv)
	}

	return srvcs
}

func createLocal() data.GruAnalytics {
	local := data.GruAnalytics{
		Service: make(map[string]data.ServiceAnalytics),
	}

	s1_is := createInstaceStatus(1, 0, 1, 1)
	s3_is := createInstaceStatus(1, 0, 0, 0)
	s4_is := createInstaceStatus(1, 0, 0, 0)
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

func createPeers() []data.GruAnalytics {
	p1 := data.GruAnalytics{
		Service: make(map[string]data.ServiceAnalytics),
	}
	p2 := data.GruAnalytics{
		Service: make(map[string]data.ServiceAnalytics),
	}

	s1_is := createInstaceStatus(1, 0, 2, 0)
	s1b_is := createInstaceStatus(1, 0, 1, 0)
	s2_is := createInstaceStatus(1, 0, 0, 1)
	s2b_is := createInstaceStatus(1, 0, 1, 0)
	s1_sa := createServiceAnalytics(0.6, 0.4, 0.4, 0.4, s1_is)
	s1b_sa := createServiceAnalytics(1.0, 0.6, 0.6, 0.6, s1b_is)
	s2_sa := createServiceAnalytics(0.4, 0.6, 0.2, 0.6, s2_is)
	s2b_sa := createServiceAnalytics(0.8, 1.0, 1.0, 0.2, s2b_is)
	s4_is := createInstaceStatus(1, 0, 0, 1)
	s4_sa := createServiceAnalytics(0.5, 0.5, 0.5, 0.5, s4_is)
	p1.Service["s1"] = s1_sa
	p1.Service["s2"] = s2b_sa
	p1.Service["s4"] = s4_sa
	p2.Service["s1"] = s1b_sa
	p2.Service["s2"] = s2_sa

	p1.System = createSystemAnalytics(
		[]string{"s1", "s2", "s4"},
		0.6,
		0.6,
		0.6,
		s1_is, s2b_is, s4_is)
	p2.System = createSystemAnalytics(
		[]string{"s1", "s2"},
		0.8,
		1.0,
		1.0,
		s1b_is, s2_is)

	p1.Health = 0.4
	p2.Health = 1.0

	return []data.GruAnalytics{p1, p2}
}

func createInstaceStatus(nInstRun int, nInstPend int, nInstStop int, nInstPaus int) cfg.ServiceStatus {
	nInstAll := nInstRun + nInstPend + nInstStop + nInstPaus

	inst_all := createRandomInstanceNames(nInstAll, 5)
	inst_run := createRandomInstanceNames(nInstRun, 5)
	inst_pen := createRandomInstanceNames(nInstPend, 5)
	inst_stp := createRandomInstanceNames(nInstStop, 5)
	inst_pau := createRandomInstanceNames(nInstPaus, 5)

	instStat := cfg.ServiceStatus{
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
	health float64, instStatus cfg.ServiceStatus) data.ServiceAnalytics {

	srvA := data.ServiceAnalytics{
		Load: load,
		Resources: data.ResourcesAnalytics{
			Cpu:    cpu,
			Memory: mem,
		},
		Instances: instStatus,
		Health:    health,
	}

	return srvA
}

func createSystemAnalytics(services []string, cpu float64,
	mem float64, health float64, instStatus ...cfg.ServiceStatus) data.SystemAnalytics {

	sysInstSt := cfg.ServiceStatus{}
	for _, st := range instStatus {
		sysInstSt.All = append(sysInstSt.All, st.All...)
		sysInstSt.Pending = append(sysInstSt.Pending, st.Pending...)
		sysInstSt.Running = append(sysInstSt.Running, st.Running...)
		sysInstSt.Stopped = append(sysInstSt.Stopped, st.Stopped...)
		sysInstSt.Paused = append(sysInstSt.Paused, st.Paused...)
	}

	sysA := data.SystemAnalytics{
		Services: services,
		Resources: data.ResourcesAnalytics{
			Cpu:    cpu,
			Memory: mem,
		},
		Instances: sysInstSt,
		Health:    health,
	}

	return sysA

}

func TestComputeLocalShared(t *testing.T) {
	analytics := data.CreateMockAnalytics()
	shared := computeLocalShared(&analytics)

	for name, _ := range analytics.Service {
		assert.Equal(t, analytics.Service[name].Load, shared.Service[name].Load)
	}

	assert.Equal(t, analytics.System.Services, shared.System.ActiveServices)
}

func TestComputeClusterData(t *testing.T) {
	defer storage.DeleteAllData(enum.SHARED)
	local := data.CreateMockShared()
	cluster := computeClusterData(local)
	assert.Equal(t, local, cluster)

	mockCluster := local
	srv := mockCluster.Service["service1"]
	srv.Load = 0.5
	mockCluster.Service["service1"] = srv
	data.SaveSharedCluster(mockCluster)
	cluster = computeClusterData(local)
	assert.NotEqual(t, local, cluster)
}

func TestAnalyzeSharedData(t *testing.T) {
	defer storage.DeleteAllData(enum.SHARED)
	var err error
	analytics := data.CreateMockAnalytics()

	analyzeSharedData(&analytics)
	_, err = data.GetSharedLocal()
	assert.NoError(t, err)
	_, err = data.GetSharedCluster()
	assert.NoError(t, err)
}
