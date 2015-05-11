package action

import (
	"github.com/samalba/dockerclient"
)

type GruActionConfig struct {
	Client      *dockerclient.DockerClient
	ContainerId string
	HostConf    *dockerclient.HostConfig
}
