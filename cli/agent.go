package cli

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/samalba/dockerclient"

	"github.com/elleFlorio/gru/action"
	"github.com/elleFlorio/gru/autonomic"
)

func agent(c *cli.Context) {

	docker, err := dockerclient.NewDockerClient("unix://var/run/docker.sock", nil)
	if err != nil {
		log.Errorln("Error in docker client creation: ", err.Error())
		return
	}

	manager := autonomic.NewAutoManager(3, docker)
	manager.RunLoop()

}
