package action

import (
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/samalba/dockerclient"
)

type GruActionConfig struct {
	Service         string
	Target          string
	TargetType      string
	Client          *dockerclient.DockerClient
	HostConfig      *dockerclient.HostConfig
	ContainerConfig *dockerclient.ContainerConfig
}
