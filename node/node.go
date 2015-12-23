package node

import (
	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/network"
	res "github.com/elleFlorio/gru/resources"
	"github.com/elleFlorio/gru/utils"
)

func CreateNode(name string, resources *res.Resource) {
	node_UUID, err := utils.GenerateUUID()
	if err != nil {
		log.WithField("err", err).Errorln("Error generating node UUID")
	}
	node_address := "http://" + network.Config().IpAddress + ":" + network.Config().Port
	config := cfg.NodeConfig{node_UUID, name, node_address, "", ""}
	nodeRes := cfg.NodeResources{
		TotalCpus:   resources.CPU.Total,
		TotalMemory: resources.Memory.Total,
	}
	node := cfg.Node{
		Configuration: config,
		Active:        false,
		Resources:     nodeRes,
	}
	cfg.SetNode(node)
}
