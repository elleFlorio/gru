package cluster

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/discovery"
	"github.com/elleFlorio/gru/node"
)

// etcd
const c_GRU_PATH = "/gru/"
const c_NODES_FOLDER = "nodes/"
const c_CONFIG_FOLDER = "config/"
const c_SERVICES_FOLDER = "services/"
const c_TTL = 5

type Cluster struct {
	UUID        string
	Name        string
	ClusterPath string
	NodePath    string
}

var (
	myCluster Cluster

	ErrNoClusterId          error = errors.New("Cannot find cluster ID")
	ErrNoCluster            error = errors.New("Node belongs to no cluster")
	ErrInvalidFriendsNumber error = errors.New("Friends number should be > 0")
	ErrInvalidDataType      error = errors.New("Invalid data type to retrieve")
	ErrNoPeers              error = errors.New("There are no peers to reach")
	ErrNoFriends            error = errors.New("There are no friends to reach")
)

func RegisterCluster(name string, id string) {
	var err error
	err = discovery.Register(c_GRU_PATH+name+"/uuid", id)
	if err != nil {
		log.Errorln("Error registering cluster")
	}
	log.Debugln("Created cluster forder: ", name)

	opt := discovery.Options{"Dir": true}
	err = discovery.Set(c_GRU_PATH+name+"/"+c_NODES_FOLDER, "", opt)
	if err != nil {
		log.Errorln("Error creating nodes folder")
	}
	log.Debugln("Created nodes folder")

	err = discovery.Set(c_GRU_PATH+name+"/"+c_SERVICES_FOLDER, "", opt)
	if err != nil {
		log.Errorln("Error creating services folder")
	}
	log.Debugln("Created services folder")

	err = discovery.Set(c_GRU_PATH+name+"/"+c_CONFIG_FOLDER, "empty", discovery.Options{})
	if err != nil {
		log.Errorln("Error creating config key")
	}
	log.Debugln("Created config key")
}

func JoinCluster(name string) error {
	key := c_GRU_PATH + name + "/uuid"
	data, err := discovery.Get(key, discovery.Options{})
	if err != nil {
		log.WithField("err", err).Errorln("Error getting cluster uuid")
		return err
	}
	if id, ok := data[key]; ok {
		initCluster(name, id, c_GRU_PATH+name, c_GRU_PATH+name+"/"+c_NODES_FOLDER+node.GetNode().Configuration.Name)
		err := createNodeFolder()
		if err != nil {
			return err
		}
		go keepAlive(c_TTL)
	} else {
		return ErrNoClusterId
	}

	return nil
}

func initCluster(name string, id string, clusterPath string, nodePath string) {
	myCluster = Cluster{
		id,
		name,
		clusterPath,
		nodePath,
	}
}

func createNodeFolder() error {
	var err error
	opt := discovery.Options{
		"TTL": time.Second * time.Duration(c_TTL),
		"Dir": true,
	}
	err = discovery.Set(myCluster.NodePath, "", opt)
	if err != nil {
		log.WithField("err", err).Errorln("Error creating node folder")
		return err
	}

	return nil
}

func keepAlive(ttl int) {
	ticker := time.NewTicker(time.Second * time.Duration(ttl))
	for {
		select {
		case <-ticker.C:
			err := updateNodeFolder(ttl)
			if err != nil {
				log.Errorln("Error keeping the node alive")
			}
		}
	}
}

func updateNodeFolder(ttl int) error {
	opt := discovery.Options{
		"TTL":       time.Second * time.Duration(ttl),
		"Dir":       true,
		"PrevExist": true,
	}
	err := discovery.Set(myCluster.NodePath, "", opt)
	WriteNodeConfig(myCluster.NodePath, node.GetNode().Configuration)
	WriteNodeActive(myCluster.NodePath, node.GetNode().Active)
	if err != nil {
		log.WithField("err", err).Errorln("Error updating node folder")
		return err
	}

	return nil
}

func GetMyCluster() (Cluster, error) {
	if myCluster.UUID == "" {
		return Cluster{}, ErrNoCluster
	}

	return myCluster, nil
}

func ListClusters() map[string]string {
	resp, err := discovery.Get(c_GRU_PATH, discovery.Options{})
	if err != nil {
		log.WithField("err", err).Errorln("Error listing clusters")
		return map[string]string{}
	}
	clusters := []string{}
	for k, _ := range resp {
		tokens := strings.Split(k, "/")
		if len(tokens) != 3 {
			// No clusters
			return map[string]string{}
		}
		clusters = append(clusters, tokens[2])
	}

	clustersUuid := make(map[string]string, len(clusters))
	for _, name := range clusters {
		resp, err := discovery.Get(c_GRU_PATH+name+"/uuid", discovery.Options{})
		if err != nil {
			log.Error("Error getting UUID of cluster ", name)
		}
		clustersUuid[name] = resp[c_GRU_PATH+name+"/uuid"]
	}

	return clustersUuid

}

func ListNodes(clusterName string, onlyActive bool) map[string]string {
	log.Debugln("Listing nodes")
	nodesPath := c_GRU_PATH + clusterName + "/nodes"

	resp, err := discovery.Get(nodesPath, discovery.Options{})
	if err != nil {
		log.WithField("err", err).Errorln("Error listing nodes in cluster ", clusterName)
		return map[string]string{}
	}

	nameAddress := make(map[string]string, len(resp))
	for nodePath, _ := range resp {
		log.Debugln("Reading configuration of node ", nodePath)
		config := &node.Config{}
		ReadNodeConfig(nodePath, config)
		if onlyActive {
			log.Debugln("Checking if node is active")
			if ReadNodeActive(nodePath) {
				nameAddress[config.Name] = config.Address
			}
		} else {
			nameAddress[config.Name] = config.Address
		}
	}

	return nameAddress
}

func GetNodes(clusterName string, onlyActive bool) []node.Node {
	nodes := []node.Node{}
	nodesPath := c_GRU_PATH + clusterName + "/nodes"
	resp, err := discovery.Get(nodesPath, discovery.Options{"Recursive": true})
	if err != nil {
		log.WithField("err", err).Errorln("Error listing nodes in cluster ", clusterName)
		return nodes
	}

	for nodePath, _ := range resp {
		n := ReadNode(nodePath)
		if onlyActive {
			if n.Active {
				nodes = append(nodes, n)
			}
		} else {
			nodes = append(nodes, n)
		}
	}

	log.Debugln("Nodes: ", nodes)

	return nodes
}

func ListServices(clusterName string) []string {
	servicesPath := c_GRU_PATH + clusterName + "/services"
	resp, err := discovery.Get(servicesPath, discovery.Options{})
	if err != nil {
		log.WithField("err", err).Errorln("Error listing services in cluster ", clusterName)
		return []string{}
	}

	services := []string{}
	for k, _ := range resp {
		path := strings.Split(k, "/")
		name := path[len(path)-1]
		services = append(services, name)
	}

	return services
}

func WriteNode(nodeData node.Node) {
	WriteNodeConfig(myCluster.NodePath, nodeData.Configuration)
	WriteNodeConstraints(myCluster.NodePath, nodeData.Constraints)
	WriteNodeResources(myCluster.NodePath, nodeData.Resources)
	WriteNodeActive(myCluster.NodePath, nodeData.Active)
}

func WriteNodeConfig(nodePath string, nodeConfig node.Config) {
	configPath := nodePath + "/config"
	err := writeData(configPath, nodeConfig)
	if err != nil {
		log.WithField("err", err).Errorln("Error writing node configuration")
	}
}

func WriteNodeConstraints(nodePath string, nodeConstraints node.Constraints) {
	constraintsPath := nodePath + "/constraints"
	err := writeData(constraintsPath, nodeConstraints)
	if err != nil {
		log.WithField("err", err).Errorln("Error writing node constraints")
	}
}

func WriteNodeResources(nodePath string, nodeResources node.Resources) {
	resourcesPath := nodePath + "/resources"
	err := writeData(resourcesPath, nodeResources)
	if err != nil {
		log.WithField("err", err).Errorln("Error writing node resources")
	}
}

func WriteNodeActive(nodePath string, active bool) {
	activePath := nodePath + "/active"
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

func writeData(path string, value interface{}) error {
	var err error
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = discovery.Set(path, string(data), discovery.Options{})
	if err != nil {
		return err
	}

	return nil
}

func ReadNode(nodePath string) node.Node {
	config := node.Config{}
	constraints := node.Constraints{}
	resources := node.Resources{}
	ReadNodeConfig(nodePath, &config)
	ReadNodeConstraints(nodePath, &constraints)
	ReadNodeResources(nodePath, &resources)
	active := ReadNodeActive(nodePath)

	return node.Node{config, constraints, resources, active}
}

func ReadNodeConfig(nodePath string, config *node.Config) {
	configPath := nodePath + "/config"
	err := readData(configPath, config)
	if err != nil {
		log.WithField("err", err).Errorln("Error reading node configuration")
	}
}

func ReadNodeConstraints(nodePath string, constraints *node.Constraints) {
	constraintsPath := nodePath + "/constraints"
	err := readData(constraintsPath, constraints)
	if err != nil {
		log.WithField("err", err).Errorln("Error reading node constraints")
	}
}

func ReadNodeResources(nodePath string, resources *node.Resources) {
	resourcesPath := nodePath + "/resources"
	err := readData(resourcesPath, resources)
	if err != nil {
		log.WithField("err", err).Errorln("Error reading node resources")
	}
}

func ReadNodeActive(nodePath string) bool {
	activePath := nodePath + "/active"
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

func readData(path string, dest interface{}) error {
	var err error
	resp, err := discovery.Get(path, discovery.Options{})
	if err != nil {
		return err
	}
	data := resp[path]
	err = json.Unmarshal([]byte(data), dest)
	if err != nil {
		return err
	}

	return nil
}
