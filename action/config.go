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
	hostConfig.CpusetCpus = sConf.CpusetCpus
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
	containerConfig := dockerclient.ContainerConfig{}
	containerConfig.Memory = getMemInBytes(sConf.Memory)
	containerConfig.Cmd = sConf.Cmd
	containerConfig.Volumes = sConf.Volumes
	containerConfig.Entrypoint = sConf.Entrypoint
	containerConfig.ExposedPorts = sConf.ExposedPorts
	containerConfig.CpuShares = sConf.CpuShares
	containerConfig.Cpuset = sConf.CpusetCpus

	return &containerConfig
}

func getMemInBytes(memory string) int64 {
	var err error
	var memInBytes int64

	if memory != "" {
		memInBytes, err = utils.RAMInBytes(memory)
		if err != nil {
			log.WithField("error", err).Warnln("Error creating container configuration.")
			return 0
		}
	} else {
		memInBytes = 0
	}

	return memInBytes
}
