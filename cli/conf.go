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
	log.WithField("status", "start").Debugln("Agent configuration loading")
	log.Debugln("Loading Agent Configuration...")

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
	log.WithField("status", "done").Debugln("Agent configuration loading")
	return &config, nil
}
