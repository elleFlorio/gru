package action

import (
	"github.com/samalba/dockerclient"
)

type GruActionConfig struct {
	Service         string
	Target          string
	TargetType      string
	Client          *dockerclient.DockerClient
	HostConfig      *dockerclient.HostConfig
	ContainerConfig *dockerclient.ContainerConfig
}
