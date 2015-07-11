package node

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

type Node struct {
	Name        string
	Constraints Constraints
}

type Constraints struct {
	CpuMin       float64
	CpuMax       float64
	MaxInstances int
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
