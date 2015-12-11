package agent

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/autonomic"
	"github.com/elleFlorio/gru/container"
	"github.com/elleFlorio/gru/metric"
	"github.com/elleFlorio/gru/storage"
)

const gruAgentConfigFile string = "/gru/config/gruagentconfig.json"

var config GruAgentConfig

func Initialize(agentConfig GruAgentConfig) {
	config = agentConfig
}

// func LoadGruAgentConfig(filename string) error {
// 	log.Infoln("Loading agent configuration")

// 	tmp, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		log.WithField("error", err).Errorln("Error reading configuration file")
// 		return err
// 	}
// 	err = json.Unmarshal(tmp, &config)
// 	if err != nil {
// 		log.WithField("error", err).Errorln("Error unmarshaling configuration file")
// 		return err
// 	}
// 	return nil
// }

func Run() {
	initializeStorage()
	initializeContainerEngine()
	initializeMetricSerivice()
	startAutonomicManager()
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

func signalErrorInAgent(err error) {
	log.WithField("err", err).Fatal("Error running gru agent. Exit.")
}

func startAutonomicManager() {
	autonomic.Initialize(
		config.Autonomic.LoopTimeInterval,
		config.Autonomic.MaxFriends)
	autonomic.RunLoop()
}

func Config() GruAgentConfig {
	return config
}
