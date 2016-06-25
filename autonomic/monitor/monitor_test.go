package monitor

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	cfg "github.com/elleFlorio/gru/configuration"
	srv "github.com/elleFlorio/gru/service"
)

func init() {
	resetMockServices()
}

func TestUpdateSystemInstances(t *testing.T) {
	defer resetMockServices()
	list := srv.List()

	updateSystemInstances(list)

	srv1, _ := srv.GetServiceByName("service1")
	srv2, _ := srv.GetServiceByName("service2")
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

func resetMockServices() {
	mockServices := srv.CreateMockServices()
	cfg.SetServices(mockServices)
}
