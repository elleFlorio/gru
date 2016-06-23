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

func WriteService(remote string, data Service) {
	err := writeData(remote, data)
	if err != nil {
		log.WithField("err", err).Errorln("Error writing service")
	}
}

func WritePolicy(remote string, data Policy) {
	err := writeData(remote, data)
	if err != nil {
		log.WithField("err", err).Errorln("Error writing policy configuration")
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
	log.WithField("remote", remote).Debugln("Reading nodes from remote")
	resp, err := discovery.Get(remote, discovery.Options{})
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
	log.WithField("remote", remote).Debugln("Reading node from remote")
	config := NodeConfig{}
	constraints := NodeConstraints{}
	resources := NodeResources{}
	ReadNodeConfig(remote, &config)
	ReadNodeConstraints(remote, &constraints)
	ReadNodeResources(remote, &resources)
	active := ReadNodeActive(remote)

	node := Node{
		Configuration: config,
		Constraints:   constraints,
		Resources:     resources,
		Active:        active,
	}

	return node
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
	resp, err := discovery.Get(remote, discovery.Options{})
	if err != nil {
		log.WithField("err", err).Errorln("Error reading services from ", remote)
		return []Service{}
	}

	services := []Service{}
	for servicePath, _ := range resp {
		srv := ReadService(servicePath)
		services = append(services, srv)
	}

	return services
}

func ReadService(remote string) Service {
	srv := Service{}
	err := readData(remote, &srv)
	if err != nil {
		log.WithField("err", err).Errorln("Error reading service")
	}

	return srv
}

func ReadAgentConfig(remote string, config *Agent) {
	err := discovery.ReadData(remote, config)
	if err != nil {
		log.WithField("err", err).Errorln("Error reading agent configuration")
	}
}

func ReadPolicyConfig(remote string) Policy {
	policy := Policy{}
	err := readData(remote, &policy)
	if err != nil {
		log.WithField("err", err).Errorln("Error reading policy parameters")
	}

	return policy
}

func ReadExpressions(remote string) map[string]Expression {
	expressions := make(map[string]Expression)
	resp, err := discovery.Get(remote, discovery.Options{})
	if err != nil {
		log.WithField("err", err).Errorln("Error reading expressions from ", remote)
		return expressions
	}

	for exprPath, _ := range resp {
		expr := ReadExpression(exprPath)
		expressions[expr.Analytic] = expr
	}

	return expressions
}

func ReadExpression(remote string) Expression {
	expr := Expression{}
	err := readData(remote, &expr)
	if err != nil {
		log.WithField("err", err).Errorln("Error reading expression")
	}

	return expr
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
