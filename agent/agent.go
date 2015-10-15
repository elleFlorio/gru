package agent

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/autonomic"
	"github.com/elleFlorio/gru/container"
	"github.com/elleFlorio/gru/discovery"
	"github.com/elleFlorio/gru/node"
	"github.com/elleFlorio/gru/service"
	"github.com/elleFlorio/gru/storage"
)

var config GruAgentConfig

func LoadGruAgentConfig(filename string) error {
	log.WithField("status", "start").Infoln("Agent configuration loading")
	defer log.WithField("status", "done").Infoln("Agent configuration loading")

	tmp, err := ioutil.ReadFile(filename)
	if err != nil {
		log.WithField("error", err).Errorln("Error reading configuration file")
		return err
	}
	err = json.Unmarshal(tmp, &config)
	if err != nil {
		log.WithField("error", err).Errorln("Error unmarshaling configuration file")
		return err
	}
	return nil
}

func Run() {
	initializeServices()
	initializeNode()
	initializeDiscovery()
	initializeStorage()
	initializeContainerEngine()

	startAutonomicManager()
}

func initializeServices() {
	servicesPath := os.Getenv("HOME") + config.Service.ServiceConfigFolder
	//Do I need to return the slice of services?
	err := service.LoadServices(servicesPath)
	if err != nil {
		signalErrorInAgent(err)
	}
}

func initializeNode() {
	nodeConfigPath := os.Getenv("HOME") + config.Node.NodeConfigFile

	err := node.LoadNodeConfig(nodeConfigPath)
	if err != nil {
		signalErrorInAgent(err)
	}
}

func initializeDiscovery() {
	log.Debugln("discovery service: ", config.Discovery.DiscoveryService)
	log.Debugln("discovery uri: ", config.Discovery.DiscoveryServiceUri)
	_, err := discovery.New(config.Discovery.DiscoveryService, config.Discovery.DiscoveryServiceUri)
	if err != nil {
		log.WithFields(log.Fields{
			"status":  "waring",
			"error":   err,
			"default": "No discovery service, running in single node mode",
		}).Warnln("Running gru agent")
	}
}

func initializeStorage() {
	_, err := storage.New(config.Storage.StorageService)
	if err != nil {
		log.WithFields(log.Fields{
			"status":  "waring",
			"error":   err,
			"default": storage.Name(),
		}).Warnln("Running gru agent")
	}
}

func initializeContainerEngine() {
	err := container.Connect(config.Docker.DaemonUrl, config.Docker.DaemonTimeout)
	if err != nil {
		signalErrorInAgent(err)
	}
	log.WithField("docker", "ok").Infoln("Container engine initialized")
}

func signalErrorInAgent(err error) {
	log.WithFields(log.Fields{
		"status": "error",
		"error":  err,
	}).Fatal("Running gru agent")
}

func startAutonomicManager() {
	autonomic.Initialize(
		config.Autonomic.LoopTimeInterval,
		config.Autonomic.MaxFrineds,
		config.Autonomic.DataToShare)
	autonomic.RunLoop()
}

func Config() GruAgentConfig {
	return config
}
