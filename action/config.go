package action

import (
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/samalba/dockerclient"

	"github.com/elleFlorio/gru/service"
	"github.com/elleFlorio/gru/utils"
)

type GruActionConfig struct {
	Service         string
	Instances       service.InstanceStatus
	HostConfig      *dockerclient.HostConfig
	ContainerConfig *dockerclient.ContainerConfig
}

func CreateHostConfig(sConf service.Config) *dockerclient.HostConfig {
	memInBytes, _ := utils.RAMInBytes(sConf.Memory)
	hostConfig := dockerclient.HostConfig{}

	hostConfig.Memory = memInBytes
	hostConfig.CpuShares = sConf.CpuShares
	hostConfig.CpusetCpus = string(sConf.CpuSet)
	hostConfig.Links = sConf.Links
	hostConfig.PortBindings = createPortBindings(sConf.PortBindings)

	return &hostConfig
}

func createPortBindings(portBindings map[string][]service.PortBinding) map[string][]dockerclient.PortBinding {
	portBindings_dckr := make(map[string][]dockerclient.PortBinding)
	for key, value := range portBindings {
		prtbndngs := []dockerclient.PortBinding{}
		for _, item := range value {
			prtbndng := dockerclient.PortBinding{}
			prtbndng.HostIp = item.HostIp
			prtbndng.HostPort = item.HostPort
			prtbndngs = append(prtbndngs, prtbndng)
		}
		portBindings_dckr[key] = prtbndngs
	}

	return portBindings_dckr
}

func CreateContainerConfig(sConf service.Config) *dockerclient.ContainerConfig {
	memInBytes, _ := utils.RAMInBytes(sConf.Memory)
	containerConfig := dockerclient.ContainerConfig{}

	containerConfig.Cmd = sConf.Cmd
	containerConfig.Volumes = sConf.Volumes
	containerConfig.Entrypoint = sConf.Entrypoint
	containerConfig.Memory = memInBytes
	containerConfig.CpuShares = sConf.CpuShares
	containerConfig.Cpuset = string(sConf.CpuSet)

	return &containerConfig
}
