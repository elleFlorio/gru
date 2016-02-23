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

const c_CMD_IP = "$ip"
const c_CMD_PORT = "$port"

var cmdSpecValue map[string]string

func init() {
	cmdSpecValue = map[string]string{
		c_CMD_IP:   "",
		c_CMD_PORT: "",
	}
}

func CreateHostConfig(name string, sConf cfg.ServiceDocker) *dockerclient.HostConfig {
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
	//hostConfig.PortBindings = createPortBindings(sConf.PortBindings)
	hostConfig.PortBindings = createPortBindings(name)

	if err != nil {
		log.Debugln("Creating Host config: Memory limit not specified")
	} else {
		hostConfig.Memory = memInBytes
	}

	return &hostConfig
}

// func createPortBindings(portBindings map[string][]cfg.PortBinding) map[string][]dockerclient.PortBinding {
// 	portBindings_dckr := make(map[string][]dockerclient.PortBinding)
// 	for key, value := range portBindings {
// 		prtbndngs := []dockerclient.PortBinding{}
// 		for _, item := range value {
// 			prtbndng := dockerclient.PortBinding{}
// 			prtbndng.HostIp = item.HostIp
// 			prtbndng.HostPort = item.HostPort
// 			prtbndngs = append(prtbndngs, prtbndng)
// 		}
// 		portBindings_dckr[key] = prtbndngs
// 	}

// 	return portBindings_dckr
// }

func createPortBindings(name string) map[string][]dockerclient.PortBinding {
	portBindings_dckr := make(map[string][]dockerclient.PortBinding)
	//TODO handle error
	res.AssignPortToService(name)
	guest := res.GetAssignedPort(name)
	host := guest
	portTcp := guest + "/tcp"
	pBindings := []dockerclient.PortBinding{}
	pBinding := dockerclient.PortBinding{
		HostIp:   "0.0.0.0",
		HostPort: host,
	}
	pBindings = append(pBindings, pBinding)

	portBindings_dckr[portTcp] = pBindings

	return portBindings_dckr
}

func CreateContainerConfig(name string, sConf cfg.ServiceDocker) *dockerclient.ContainerConfig {
	containerConfig := dockerclient.ContainerConfig{}
	containerConfig.Memory = getMemInBytes(sConf.Memory)
	containerConfig.Env = getEnvVars(sConf.Env)
	containerConfig.Volumes = sConf.Volumes
	containerConfig.Entrypoint = sConf.Entrypoint
	//containerConfig.ExposedPorts = sConf.ExposedPorts
	containerConfig.ExposedPorts = getExposedPorts(name)
	containerConfig.CpuShares = sConf.CpuShares
	containerConfig.Cpuset = sConf.CpusetCpus
	//containerConfig.Cmd = sConf.Cmd
	containerConfig.Cmd = getCommands(sConf.Cmd)

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

func getExposedPorts(name string) map[string]struct{} {
	exposed := make(map[string]struct{})
	guest := res.GetAssignedPort(name)
	guestTcp := guest + "/tcp"
	exposed[guestTcp] = struct{}{}

	return exposed
}

func getCommands(cmds map[string]string) []string {
	commands := make([]string, 0, len(cmds))
	for cmd, value := range cmds {
		switch value {
		case c_CMD_IP:
			value = cmdSpecValue[c_CMD_IP]
		case c_CMD_PORT:
			value = cmdSpecValue[c_CMD_PORT]
		}

		commands = append(commands, cmd, value)
	}

	return commands
}
