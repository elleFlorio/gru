package node

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/network"
	"github.com/elleFlorio/gru/utils"
)

type Node struct {
	UUID        string      `json:"uuid"`
	Name        string      `json:"name"`
	IpAddr      string      `json:"ipaddr"`
	Port        string      `json:"port"`
	Constraints Constraints `json:"constraints"`
}

type Constraints struct {
	CpuMin       float64 `json:"cpumin"`
	CpuMax       float64 `json:"cpumax"`
	MaxInstances int     `json:"maxinstances"`
}

var node Node

func LoadNodeConfig(filename string) error {
	node.UUID, _ = utils.GenerateUUID()

	log.WithField("status", "start").Infoln("Node configuration loading")
	defer log.WithFields(log.Fields{
		"status": "done",
		"UUID":   node.UUID,
	}).Infoln("Node configuration loading")

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

	if node.IpAddr == "" {
		node.IpAddr, err = network.GetHostIp()
		if err != nil {
			log.WithField("error", err).Errorln("Error retrieving ip address")
			return err
		}
	}

	if node.Port == "" {
		node.Port, err = network.GetPort()
		if err != nil {
			log.WithField("error", err).Errorln("Error retrieving port. Set to default [5000]")
		}
	}

	log.WithFields(log.Fields{
		"ipAddr": node.IpAddr,
		"port":   node.Port,
	}).Infoln("Node address set")

	return nil
}

func GetNodeConfig() Node {
	return node
}

func UpdateNode(newNode Node) {
	node = newNode
}
