package action

import (
	"encoding/json"
	"os"

	log "github.com/Sirupsen/logrus"
)

type GruActionConfig struct {
	Services struct {
		MinActive int
		MaxActive int
		Service   map[string]int
	}

	Nodes struct {
		MaxInstances int
		Node         map[string]int
	}

	Cpu struct {
		Min     float64
		Max     float64
		Service map[string][2]float64
	}

	Docker struct {
		StopTimeout int
	}
}

func BuildGruActionConfig(filePath string) (*GruActionConfig, error) {
	var config GruActionConfig
	configFile, err := os.Open(filePath)
	if err != nil {
		log.Errorln("Error opening action config file", err.Error())
		return &config, err
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		log.Errorln("Error parsing action config file", err.Error())
		return &config, err
	}

	log.Infoln("Gru Action Configuration built")

	return &config, nil
}
