package action

import (
	"github.com/samalba/dockerclient"
)

type GruActionConfig struct {
	Service    string
	Target     string
	Client     *dockerclient.DockerClient
	HostConfig *dockerclient.HostConfig
}
