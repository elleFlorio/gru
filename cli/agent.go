package cli

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/samalba/dockerclient"

	"github.com/elleFlorio/gru/autonomic"
	"github.com/elleFlorio/gru/service"
)

const gruAgentConfigFolder string = "/gru/config/gruagentconfig.json"
const servicesFolder string = "/gru/config/services"

func agent(c *cli.Context) {
	log.WithField("status", "start").Infoln("Running gru agent")
	defer log.WithField("status", "done").Infoln("Running gru agent")

	gruAgentConfigPath := os.Getenv("HOME") + gruAgentConfigFolder
	servicesPath := os.Getenv("HOME") + servicesFolder

	config, err := LoadGruAgentConfig(gruAgentConfigPath)
	if err != nil {
		signalErrorInAgent(err)
		return
	}

	//Do I need to return the slice of services?
	_, err = service.LoadServices(servicesPath)
	if err != nil {
		signalErrorInAgent(err)
		return
	}

	docker, err := dockerclient.NewDockerClient(config.DaemonUrl, nil)
	if err != nil {
		signalErrorInAgent(err)
		return
	}

	log.WithField("status", "config complete").Infoln("Running gru agent")

	manager := autonomic.NewAutoManager(docker, config.LoopTimeInterval)
	manager.RunLoop()
}

func signalErrorInAgent(err error) {
	log.WithFields(log.Fields{
		"status": "error",
		"error":  err,
	}).Errorln("Running gru agent")

	log.WithField("status", "abort").Errorln("Running gru agent")
}
