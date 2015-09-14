package cli

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/samalba/dockerclient"

	"github.com/elleFlorio/gru/api"
	"github.com/elleFlorio/gru/autonomic"
	"github.com/elleFlorio/gru/discovery"
	"github.com/elleFlorio/gru/node"
	"github.com/elleFlorio/gru/service"
	"github.com/elleFlorio/gru/storage"
)

const gruAgentConfigFile string = "/gru/config/gruagentconfig.json"
const servicesFolder string = "/gru/config/services"
const nodeConfigFile string = "/gru/config/nodeconfig.json"

func agent(c *cli.Context) {
	log.WithField("status", "start").Infoln("Running gru agent")
	defer log.WithField("status", "done").Infoln("Running gru agent")

	gruAgentConfigPath := os.Getenv("HOME") + gruAgentConfigFile
	servicesPath := os.Getenv("HOME") + servicesFolder
	nodeConfigPath := os.Getenv("HOME") + nodeConfigFile

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

	err = node.LoadNodeConfig(nodeConfigPath)
	if err != nil {
		signalErrorInAgent(err)
		return
	}

	docker, err := dockerclient.NewDockerClient(config.DaemonUrl, nil)
	if err != nil {
		signalErrorInAgent(err)
		return
	}

	_, err = discovery.New(config.DiscoveryService, config.DiscoveryServiceUri)
	if err != nil {
		log.WithFields(log.Fields{
			"status":  "waring",
			"error":   err,
			"default": "No discovery service, running in single node mode",
		}).Warnln("Running gru agent")
	}

	//TODO this should be parametrized in configuration
	_, err = storage.New("internal")
	if err != nil {
		log.WithFields(log.Fields{
			"status":  "waring",
			"error":   err,
			"default": storage.DataStore().Name(),
		}).Warnln("Running gru agent")
	}
	storage.DataStore().Initialize()

	agentPort := ":" + node.Config().Port
	go api.StartServer(agentPort)

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
