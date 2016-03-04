package executor

import (
	"os"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/samalba/dockerclient"

	"github.com/elleFlorio/gru/autonomic/executor/action"
	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/enum"
	res "github.com/elleFlorio/gru/resources"
	"github.com/elleFlorio/gru/utils"
)

func buildConfig(srv *cfg.Service, act enum.Action) action.Action {
	actConfig := action.Action{}
	actConfig.HostConfig = createHostConfig(srv, act)
	actConfig.ContainerConfig = createContainerConfig(srv, act)
	actConfig.Parameters.DiscoveryPort = getDiscoveryPort(srv, act)
	actConfig.Service = srv.Name
	actConfig.Instances = srv.Instances
	actConfig.ContainerConfig.Image = srv.Image
	actConfig.Parameters.StopTimeout = srv.Docker.StopTimeout

	return actConfig
}

func createHostConfig(srv *cfg.Service, act enum.Action) *dockerclient.HostConfig {
	hostConfig := dockerclient.HostConfig{}

	if act == enum.START {
		conf := srv.Docker

		hostConfig.CpusetCpus = createCpusetCpus(conf.CpusetCpus, conf.CPUnumber)
		hostConfig.Memory = createMemory(conf.Memory)
		hostConfig.PortBindings = createPortBindings(srv.Name)
		hostConfig.CpuShares = conf.CpuShares
		hostConfig.Links = conf.Links
	}

	return &hostConfig
}

func createCpusetCpus(cpusetcpus string, cores int) string {
	if cpusetcpus == "" {
		if cores < 1 {
			log.Warnln("Number of requested CPUs = 0. Setting to 1")
			cores = 1
		}
		if assigned, ok := res.GetCoresAvailable(cores); ok {
			cpusetcpus = assigned
		} else {
			log.Debugln("Error setting cpusetcpus in hostconfig")
		}
	}

	return cpusetcpus
}

func createMemory(memory string) int64 {
	memInBytes, err := utils.RAMInBytes(memory)
	if err != nil {
		log.Debugln("Creating Host config: Memory limit not specified")
		return 0
	}

	return memInBytes
}

func createPortBindings(name string) map[string][]dockerclient.PortBinding {
	portBindings_dckr := make(map[string][]dockerclient.PortBinding)

	assigned, err := res.RequestPortsForService(name)
	if err != nil {
		log.Errorln("Error creating port bindings")
		return portBindings_dckr
	}
	for guest, host := range assigned {
		portTcp := guest + "/tcp"
		pBindings := []dockerclient.PortBinding{}
		pBinding := dockerclient.PortBinding{
			HostIp:   "0.0.0.0",
			HostPort: host,
		}
		pBindings = append(pBindings, pBinding)
		portBindings_dckr[portTcp] = pBindings
	}

	return portBindings_dckr
}

func createContainerConfig(srv *cfg.Service, act enum.Action) *dockerclient.ContainerConfig {
	containerConfig := dockerclient.ContainerConfig{}

	if act == enum.START {
		conf := srv.Docker
		containerConfig.Memory = createMemory(conf.Memory)
		containerConfig.Env = createEnvVars(conf.Env)
		containerConfig.ExposedPorts = createExposedPorts(srv.Name)
		containerConfig.Cmd = conf.Cmd
		containerConfig.Volumes = conf.Volumes
		containerConfig.Entrypoint = conf.Entrypoint
		containerConfig.CpuShares = conf.CpuShares
		containerConfig.Cpuset = conf.CpusetCpus
	}

	return &containerConfig
}

func createEnvVars(vars map[string]string) []string {
	envVars := make([]string, 0, len(vars))
	for name, value := range vars {
		if value == "" {
			value = os.Getenv(name)
			if value == "" {
				log.WithField("var", name).Warnln("Cannot get value of env variable")
			}
		}
		envVars = append(envVars, name+"="+value)
	}

	return envVars
}

func createExposedPorts(name string) map[string]struct{} {
	exposed := make(map[string]struct{})
	assigned := res.GetAssignedPorts(name)
	for guest, _ := range assigned {
		guestTcp := guest + "/tcp"
		exposed[guestTcp] = struct{}{}
	}

	return exposed
}

// The discovery port is the first set by the user
// in port bindings.
func getDiscoveryPort(srv *cfg.Service, act enum.Action) string {
	if act == enum.START {
		assigned := res.GetResources().Network.ServicePorts[srv.Name].LastRequested
		for _, host := range assigned {
			return host
		}
	}

	return ""
}
