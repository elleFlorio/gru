package node

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

type Node struct {
	Name        string      `json:"name"`
	Constraints Constraints `json:"constraints"`
}

type Constraints struct {
	CpuMin       float64 `json:"cpumin"`
	CpuMax       float64 `json:"cpumax"`
	MaxInstances int     `json:"maxinstances"`
}

var node Node

func LoadNodeConfig(filename string) error {
	log.WithField("status", "start").Infoln("Node configuration loading")
	defer log.WithField("status", "done").Infoln("Node configuration loading")
	tmp, err := ioutil.ReadFile(filename)
	if err != nil {
		log.WithField("error", err).Errorln("Error reading node configuration file")
		return err
	}
	err = json.Unmarshal(tmp, &node)
	if err != nil {
		log.WithField("error", err).Errorln("Error unmarshaling node configuration file")
		return err
	}

	return nil
}

func GetNodeConfig() *Node {
	return &node
}

func UpdateNode(newNode Node) {
	node = newNode
}
