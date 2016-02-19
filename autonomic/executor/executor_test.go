package executor

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/autonomic/executor/action"
	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/resources"
	"github.com/elleFlorio/gru/service"
)

func TestGetTargetService(t *testing.T) {
	defer cfg.CleanServices()
	services := service.CreateMockServices()
	cfg.SetServices(services)

	var srv *cfg.Service
	srv = getTargetService("service1")
	assert.Equal(t, "service1", srv.Name)

	srv = getTargetService("noservice")
	assert.Equal(t, "noservice", srv.Name)

	srv = getTargetService("pippo")
	assert.Equal(t, "noservice", srv.Name)
}

func TestBuildConfig(t *testing.T) {
	defer cfg.CleanServices()
	services := service.CreateMockServices()
	cfg.SetServices(services)
	resources.CreateMockResources(2, "1G", 0, "0G")

	var srv *cfg.Service
	var conf action.GruActionConfig
	srv, _ = service.GetServiceByName("service1")
	conf = buildConfig(srv)
	assert.NotEmpty(t, conf)
	srv = &cfg.Service{Name: "noservice"}
	conf = buildConfig(srv)
	assert.NotEmpty(t, conf)
}
