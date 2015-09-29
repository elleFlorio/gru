package container

import (
	"time"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/samalba/dockerclient"
)

type DockerConfig struct {
	DaemonUrl     string
	DaemonTimeout int
	Client        *dockerclient.DockerClient
}

var docker DockerConfig

//TODO implement tls config
func Connect(daemonUrl string, timeout int) error {
	client, err := dockerclient.NewDockerClientTimeout(daemonUrl, nil, time.Duration(timeout)*time.Second)
	if err != nil {
		return err
	}

	docker.DaemonUrl = daemonUrl
	docker.DaemonTimeout = timeout
	docker.Client = client

	return nil
}

func Docker() DockerConfig {
	return docker
}
