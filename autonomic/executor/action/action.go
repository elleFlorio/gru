package action

import (
	"os"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/samalba/dockerclient"

	cfg "github.com/elleFlorio/gru/configuration"
	res "github.com/elleFlorio/gru/resources"
	"github.com/elleFlorio/gru/utils"
)

type GruActionConfig struct {
	Service         string
	Instances       cfg.ServiceStatus
	Parameters      ActionParameters
	HostConfig      *dockerclient.HostConfig
	ContainerConfig *dockerclient.ContainerConfig
}

type ActionParameters struct {
	StopTimeout int
}

func CreateHostConfig(sConf cfg.ServiceDocker) *dockerclient.HostConfig {
	memInBytes, err := utils.RAMInBytes(sConf.Memory)
	hostConfig := dockerclient.HostConfig{}

	hostConfig.CpuShares = sConf.CpuShares
	if sConf.CpusetCpus == "" {
		if assigned, ok := res.GetCoresAvailable(sConf.CPUnumber); ok {
			hostConfig.CpusetCpus = assigned
		} else {
			log.Errorln("Error setting cpusetcpus in hostconfig")
		}
	}
	hostConfig.Links = sConf.Links
	hostConfig.PortBindings = createPortBindings(sConf.PortBindings)

	if err != nil {
		log.Debugln("Creating Host config: Memory limit not specified")
	} else {
		hostConfig.Memory = memInBytes
	}

	return &hostConfig
}

func createPortBindings(portBindings map[string][]cfg.PortBinding) map[string][]dockerclient.PortBinding {
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

func CreateContainerConfig(sConf cfg.ServiceDocker) *dockerclient.ContainerConfig {
	containerConfig := dockerclient.ContainerConfig{}
	containerConfig.Memory = getMemInBytes(sConf.Memory)
	containerConfig.Env = getEnvVars(sConf.Env)
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
			log.WithField("err", err).Warnln("Error creating container configuration")
			return 0
		}
	} else {
		memInBytes = 0
	}

	return memInBytes
}

func getEnvVars(vars map[string]string) []string {
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
