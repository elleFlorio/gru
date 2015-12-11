package cli

import (
	"os"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/codegangsta/cli"

	"github.com/elleFlorio/gru/agent"
	"github.com/elleFlorio/gru/api"
	"github.com/elleFlorio/gru/network"
)

const gruAgentConfigFile string = "/gru/config/gruagentconfig.json"

func start(c *cli.Context) {
	log.Infoln("Starting gru agent")
	defer agent.Run()

	gruAgentConfigPath := os.Getenv("HOME") + gruAgentConfigFile

	err := agent.LoadGruAgentConfig(gruAgentConfigPath)
	if err != nil {
		log.WithFields(log.Fields{
			"status": "error",
			"error":  err,
		}).Fatalln("Starting gru agent")
	}

	network.InitializeNetwork(agent.Config().Network.IpAddres, agent.Config().Network.Port)
	go api.StartServer(network.Config().Port)

	log.WithField("status", "ok").Infoln("Gru agent started")
}
