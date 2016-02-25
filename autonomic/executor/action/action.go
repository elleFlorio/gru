package action

import (
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/samalba/dockerclient"

	cfg "github.com/elleFlorio/gru/configuration"
)

type Action struct {
	Service         string
	Instances       cfg.ServiceStatus
	Parameters      ActionParameters
	HostConfig      *dockerclient.HostConfig
	ContainerConfig *dockerclient.ContainerConfig
}

type ActionParameters struct {
	StopTimeout   int
	DiscoveryPort string
}
