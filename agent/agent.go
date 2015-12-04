package agent

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/autonomic"
	"github.com/elleFlorio/gru/container"
	"github.com/elleFlorio/gru/discovery"
	"github.com/elleFlorio/gru/metric"
	"github.com/elleFlorio/gru/node"
	"github.com/elleFlorio/gru/service"
	"github.com/elleFlorio/gru/storage"
)

var config GruAgentConfig

func LoadGruAgentConfig(filename string) error {
	log.Infoln("Loading agent configuration")

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
	initializeDiscovery()
	initializeStorage()
	initializeContainerEngine()
	initializeMetricSerivice()
	initializeServices()
	initializeNode()

	startAutonomicManager()
}

func initializeDiscovery() {
	_, err := discovery.New(config.Discovery.DiscoveryService, config.Discovery.DiscoveryServiceUri)
	if err != nil {
		log.WithFields(log.Fields{
			"status":  "warning",
			"error":   err,
			"default": "No discovery service, running in single node mode",
		}).Warnln("Error initializing discovery service")
	} else {
		log.WithField(discovery.Name(), "ok").Infoln("Discovery service initialized")
	}
}

func initializeStorage() {
	_, err := storage.New(config.Storage.StorageService)
	if err != nil {
		log.WithFields(log.Fields{
			"status":  "warning",
			"error":   err,
			"default": storage.Name(),
		}).Warnln("Error initializing storage service")
	} else {
		log.WithField(storage.Name(), "ok").Infoln("Storage service initialized")
	}
}

func initializeContainerEngine() {
	err := container.Connect(config.Docker.DaemonUrl, config.Docker.DaemonTimeout)
	if err != nil {
		signalErrorInAgent(err)
	}
	log.WithField("docker", "ok").Infoln("Container engine initialized")
}

func initializeMetricSerivice() {
	_, err := metric.New(config.Metric.MetricService, config.Metric.Configuration)
	if err != nil {
		log.WithFields(log.Fields{
			"status":  "warning",
			"error":   err,
			"default": metric.Name(),
		}).Warnln("Error initializing metric service")
	} else {
		log.WithField(metric.Name(), "ok").Infoln("Metric service initialized")
	}
}

func initializeServices() {
	servicesPath := os.Getenv("HOME") + config.Service.ServiceConfigFolder
	//Do I need to return the slice of services?
	err := service.LoadServices(servicesPath)
	if err != nil {
		signalErrorInAgent(err)
	}
	log.WithFields(log.Fields{
		"services": "ok",
		"loaded":   len(service.List()),
	}).Infoln("Services initialized")
}

func initializeNode() {
	nodeConfigPath := os.Getenv("HOME") + config.Node.NodeConfigFile

	err := node.LoadNodeConfig(nodeConfigPath)
	if err != nil {
		signalErrorInAgent(err)
	}
	node.ComputeTotalResources()
	log.WithFields(log.Fields{
		"node": "ok",
		"UUID": node.Config().UUID,
	}).Infoln("Node initialized")
}

func signalErrorInAgent(err error) {
	log.WithField("err", err).Fatal("Error running gru agent. Exit.")
}

func startAutonomicManager() {
	autonomic.Initialize(
		config.Autonomic.LoopTimeInterval,
		config.Autonomic.MaxFriends,
		config.Autonomic.DataToShare)
	autonomic.RunLoop()
}

func Config() GruAgentConfig {
	return config
}
