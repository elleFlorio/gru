package executor

import (
	"os"
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	//"github.com/elleFlorio/gru/autonomic/executor/action"
	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/enum"
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

func TestCreateCpusetCpus(t *testing.T) {
	resources.CreateMockResources(2, "1G", 0, "0G")
	empty := ""
	core0 := "0"
	var cpusetcpus string

	cpusetcpus = createCpusetCpus(empty, 1)
	assert.Equal(t, "0", cpusetcpus)

	cpusetcpus = createCpusetCpus(core0, 1)
	assert.Equal(t, "0", cpusetcpus)

	cpusetcpus = createCpusetCpus(empty, 0)
	assert.Equal(t, "0", cpusetcpus)
}

func TestCreateMemory(t *testing.T) {
	notSet := ""
	mem1 := "1G"
	var memory int64

	memory = createMemory(notSet)
	assert.Equal(t, int64(0), memory)

	memory = createMemory(mem1)
	assert.Equal(t, 1024, 1024)
}

func TestCreatePortBindings(t *testing.T) {
	service1 := "service1"
	ports1 := map[string]string{
		"50100": "50100",
	}
	service2 := "service2"
	ports2 := map[string]string{
		"50200": "50200",
	}
	service3 := "service3"
	resources.InitializeServiceAvailablePorts(service1, ports1)
	resources.InitializeServiceAvailablePorts(service2, ports2)

	bindings1 := createPortBindings(service1)
	assert.Len(t, bindings1, 1)
	assert.NotEmpty(t, bindings1["50100/tcp"])
	assert.Equal(t, "50100", bindings1["50100/tcp"][0].HostPort)

	bindings2 := createPortBindings(service2)
	assert.Len(t, bindings2, 1)
	assert.NotEmpty(t, bindings2["50200/tcp"])
	assert.Equal(t, "50200", bindings2["50200/tcp"][0].HostPort)

	bindings3 := createPortBindings(service3)
	assert.Empty(t, bindings3)
}

func TestCreateHostConfig(t *testing.T) {
	defer cfg.CleanServices()
	services := service.CreateMockServices()
	cfg.SetServices(services)
	resources.CreateMockResources(2, "1G", 0, "0G")
	ports1 := map[string]string{
		"50100": "50100",
	}
	resources.InitializeServiceAvailablePorts("service1", ports1)
	service1, _ := service.GetServiceByName("service1")

	hostConfigStop := createHostConfig(service1, enum.STOP)
	assert.Equal(t, hostConfigStop.CpusetCpus, "")

	hostConfigStart := createHostConfig(service1, enum.START)
	assert.Equal(t, hostConfigStart.CpusetCpus, "0")
	assert.Len(t, hostConfigStart.PortBindings, 1)
}

func TestCreateEnvVars(t *testing.T) {
	vars := map[string]string{
		"pippo":    "topolinia",
		"paperino": "paperopoli",
		"topolino": "",
		"etabeta":  "",
	}
	os.Setenv("topolino", "topolinia")

	envVars := createEnvVars(vars)
	assert.Contains(t, envVars, "pippo=topolinia")
	assert.Contains(t, envVars, "paperino=paperopoli")
	assert.Contains(t, envVars, "topolino=topolinia")
	assert.Contains(t, envVars, "etabeta=")
}

func TestCreateExposedPorts(t *testing.T) {
	service1 := "service1"
	id1 := "pippo"
	availablePorts1 := map[string]string{
		"50100": "50100",
	}
	ports1 := map[string][]string{
		"50100": []string{"50100"},
	}

	resources.InitializeServiceAvailablePorts(service1, availablePorts1)
	resources.AssignSpecifiPortsToService(service1, id1, ports1)

	exposed := createExposedPorts(service1)
	assert.Len(t, exposed, 1)
	assert.Equal(t, exposed["50100/tcp"], struct{}{})
}

func TestCreateContainerConfig(t *testing.T) {
	defer cfg.CleanServices()
	services := service.CreateMockServices()
	cfg.SetServices(services)
	resources.CreateMockResources(2, "1G", 0, "0G")
	service1 := "service1"
	id1 := "pippo"
	availablePorts1 := map[string]string{
		"50100": "50100",
	}
	ports1 := map[string][]string{
		"50100": []string{"50100"},
	}
	resources.InitializeServiceAvailablePorts(service1, availablePorts1)
	resources.AssignSpecifiPortsToService(service1, id1, ports1)
	srv1, _ := service.GetServiceByName(service1)

	containerConfigStop := createContainerConfig(srv1, enum.STOP)
	assert.Empty(t, containerConfigStop.ExposedPorts)
	containerConfigStart := createContainerConfig(srv1, enum.START)
	assert.NotEmpty(t, containerConfigStart.ExposedPorts)
}

func TestBuildConfig(t *testing.T) {
	defer cfg.CleanServices()
	services := service.CreateMockServices()
	cfg.SetServices(services)
	resources.CreateMockResources(2, "1G", 0, "0G")
	ports1 := map[string]string{
		"50100": "50100",
	}
	resources.InitializeServiceAvailablePorts("service1", ports1)
	service1, _ := service.GetServiceByName("service1")

	config := buildConfig(service1, enum.START)
	assert.Equal(t, "0", config.HostConfig.CpusetCpus)
}
