package configuration

import (
	"encoding/json"
	"strconv"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/discovery"
)

func WriteNode(remote string, data Node) {
	WriteNodeConfig(remote, data.Configuration)
	WriteNodeConstraints(remote, data.Constraints)
	WriteNodeResources(remote, data.Resources)
	WriteNodeActive(remote, data.Active)
}

func WriteNodeConfig(remote string, data NodeConfig) {
	configPath := remote + "/config"
	err := writeData(configPath, data)
	if err != nil {
		log.WithField("err", err).Errorln("Error writing node configuration")
	}
}

func WriteNodeConstraints(remote string, data NodeConstraints) {
	constraintsPath := remote + "/constraints"
	err := writeData(constraintsPath, data)
	if err != nil {
		log.WithField("err", err).Errorln("Error writing node constraints")
	}
}

func WriteNodeResources(remote string, data NodeResources) {
	resourcesPath := remote + "/resources"
	err := writeData(resourcesPath, data)
	if err != nil {
		log.WithField("err", err).Errorln("Error writing node resources")
	}
}

func WriteNodeActive(remote string, active bool) {
	activePath := remote + "/active"
	var value string
	if active {
		value = "true"
	} else {
		value = "false"
	}
	err := discovery.Set(activePath, value, discovery.Options{})
	if err != nil {
		log.WithField("err", err).Errorln("Error writing node active")
	}

}

func writeData(remote string, src interface{}) error {
	var err error
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}

	err = discovery.Set(remote, string(data), discovery.Options{})
	if err != nil {
		return err
	}

	return nil
}

func ReadNodes(remote string) []Node {
	resp, err := discovery.Get(remote, discovery.Options{"Recursive": true})
	if err != nil {
		log.WithField("err", err).Errorln("Error reading nodes from ", remote)
		return []Node{}
	}

	nodes := []Node{}
	for nodePath, _ := range resp {
		n := ReadNode(nodePath)
		nodes = append(nodes, n)
	}

	return nodes
}

func ReadNode(remote string) Node {
	config := NodeConfig{}
	constraints := NodeConstraints{}
	resources := NodeResources{}
	ReadNodeConfig(remote, &config)
	ReadNodeConstraints(remote, &constraints)
	ReadNodeResources(remote, &resources)
	active := ReadNodeActive(remote)

	return Node{config, constraints, resources, active}
}

func ReadNodeConfig(remote string, config *NodeConfig) {
	configPath := remote + "/config"
	err := readData(configPath, config)
	if err != nil {
		log.WithField("err", err).Errorln("Error reading node configuration")
	}
}

func ReadNodeConstraints(remote string, constraints *NodeConstraints) {
	constraintsPath := remote + "/constraints"
	err := readData(constraintsPath, constraints)
	if err != nil {
		log.WithField("err", err).Errorln("Error reading node constraints")
	}
}

func ReadNodeResources(remote string, resources *NodeResources) {
	resourcesPath := remote + "/resources"
	err := readData(resourcesPath, resources)
	if err != nil {
		log.WithField("err", err).Errorln("Error reading node resources")
	}
}

func ReadNodeActive(remote string) bool {
	activePath := remote + "/active"
	resp, err := discovery.Get(activePath, discovery.Options{})
	if err != nil {
		log.WithField("err", err).Errorln("Error reading node active")
		return false
	}
	active, err := strconv.ParseBool(resp[activePath])
	if err != nil {
		log.WithField("err", err).Errorln("Error parsing node active")
		return false
	}

	return active
}

func ReadServices(remote string) []Service {
	resp, err := discovery.Get(remote, discovery.Options{"Recursive": true})
	if err != nil {
		log.WithField("err", err).Errorln("Error reading services from ", remote)
		return []Service{}
	}

	services := []Service{}
	for servicePath, _ := range resp {
		srv := Service{}
		ReadService(servicePath, &srv)
		services = append(services, srv)
	}

	return services
}

func ReadService(remote string, srv *Service) {
	err := readData(remote, srv)
	if err != nil {
		log.WithField("err", err).Errorln("Error reading service")
	}
}

func ReadAgentConfig(remote string, config *Agent) {
	err := discovery.ReadData(remote, config)
	if err != nil {
		log.WithField("err", err).Errorln("Error reading agent configuration")
	}
}

func readData(remote string, dest interface{}) error {
	var err error
	resp, err := discovery.Get(remote, discovery.Options{})
	if err != nil {
		return err
	}
	data := resp[remote]
	err = json.Unmarshal([]byte(data), dest)
	if err != nil {
		return err
	}

	return nil
}
