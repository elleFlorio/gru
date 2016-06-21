package monitor

import (
	"testing"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

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
