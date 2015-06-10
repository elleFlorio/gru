package cli

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

type GruAgentConfig struct {
	DaemonUrl           string
	DaemonTimeout       int
	LoopTimeInterval    int
	ServiceConfigFolder string
}

var config GruAgentConfig

func LoadGruAgentConfig(filename string) (*GruAgentConfig, error) {
	log.WithField("status", "start").Infoln("Agent configuration loading")
	defer log.WithField("status", "done").Infoln("Agent configuration loading")

	tmp, err := ioutil.ReadFile(filename)
	if err != nil {
		log.WithField("error", err).Errorln("Error reading configuration file")
		return nil, err
	}
	err = json.Unmarshal(tmp, &config)
	if err != nil {
		log.WithField("error", err).Errorln("Error unmarshaling configuration file")
		return nil, err
	}
	return &config, nil
}
