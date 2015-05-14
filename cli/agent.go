package cli

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/samalba/dockerclient"

	"github.com/elleFlorio/gru/autonomic"
	"github.com/elleFlorio/gru/service"
)

const servicesPath string = "config/services"
const gruAgentConfigPath string = "config/gruagentconfig.json"

func agent(c *cli.Context) {
	log.WithField("status", "start").Debugln("Running gru agent")
	defer log.WithField("status", "done").Errorln("Running gru agent")

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

	log.WithField("status", "config complete").Debugln("Running gru agent")

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
