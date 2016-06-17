package container

import (
	"strings"
	"time"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
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

	client, err := dockerclient.NewDockerClientTimeout(daemonUrl, nil, time.Duration(timeout)*time.Second, nil)
	if err != nil {
		return err
	}

	if _, err = client.Info(); err != nil {
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

func GetPortBindings(id string) (map[string][]string, error) {
	info, err := docker.Client.InspectContainer(id)
	if err != nil {
		log.WithFields(log.Fields{
			"id":  id,
			"err": err,
		}).Errorln("Error inspecting instance")
	}

	portBindings := createPortBindings(info.HostConfig.PortBindings)

	return portBindings, err
}

func createPortBindings(dockerBindings map[string][]dockerclient.PortBinding) map[string][]string {
	portBindings := make(map[string][]string)

	for guestTcp, bindings := range dockerBindings {
		guest := strings.Split(guestTcp, "/")[0]
		hosts := make([]string, 0, len(bindings))
		for _, host := range bindings {
			hosts = append(hosts, host.HostPort)
		}
		portBindings[guest] = hosts
	}

	return portBindings
}
