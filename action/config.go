package action

import (
	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
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
	memInBytes, err := utils.RAMInBytes(sConf.Memory)
	hostConfig := dockerclient.HostConfig{}

	hostConfig.CpuShares = sConf.CpuShares
	hostConfig.CpusetCpus = string(sConf.CpuSet)
	hostConfig.Links = sConf.Links
	hostConfig.PortBindings = createPortBindings(sConf.PortBindings)

	if err != nil {
		log.Warnln("Creating Host config: Memory limit not specified")
	} else {
		hostConfig.Memory = memInBytes
	}

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
	memInBytes, err := utils.RAMInBytes(sConf.Memory)
	containerConfig := dockerclient.ContainerConfig{}

	containerConfig.Cmd = sConf.Cmd
	containerConfig.Volumes = sConf.Volumes
	containerConfig.Entrypoint = sConf.Entrypoint
	containerConfig.CpuShares = sConf.CpuShares
	containerConfig.Cpuset = string(sConf.CpuSet)

	if err != nil {
		log.Warnln("Creating Container config: Memory limit not specified")
	} else {
		containerConfig.Memory = memInBytes
	}

	return &containerConfig
}
