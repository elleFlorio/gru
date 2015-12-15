package cluster

import (
	"errors"
	"strings"
	"time"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/discovery"
)

// etcd
const c_GRU_REMOTE = "/gru/"
const c_NODES_REMOTE = "nodes/"
const c_CONFIG_REMOTE = "config/"
const c_SERVICES_REMOTE = "services/"
const c_TTL = 5

type Cluster struct {
	UUID   string
	Name   string
	Remote string
}

var (
	myCluster Cluster

	ErrNoClusterId error = errors.New("Cannot find cluster ID")
	ErrNoCluster   error = errors.New("Node belongs to no cluster")
)

func RegisterCluster(name string, id string) {
	var err error
	err = discovery.Register(c_GRU_REMOTE+name+"/uuid", id)
	if err != nil {
		log.Errorln("Error registering cluster")
	}
	log.Debugln("Created cluster forder: ", name)

	opt := discovery.Options{"Dir": true}
	err = discovery.Set(c_GRU_REMOTE+name+"/"+c_NODES_REMOTE, "", opt)
	if err != nil {
		log.Errorln("Error creating nodes folder")
	}
	log.Debugln("Created nodes folder")

	err = discovery.Set(c_GRU_REMOTE+name+"/"+c_SERVICES_REMOTE, "", opt)
	if err != nil {
		log.Errorln("Error creating services folder")
	}
	log.Debugln("Created services folder")

	err = discovery.Set(c_GRU_REMOTE+name+"/"+c_CONFIG_REMOTE, "empty", discovery.Options{})
	if err != nil {
		log.Errorln("Error creating config key")
	}
	log.Debugln("Created config key")
}

func JoinCluster(name string) error {
	key := c_GRU_REMOTE + name + "/uuid"
	data, err := discovery.Get(key, discovery.Options{})
	if err != nil {
		log.WithField("err", err).Errorln("Error getting cluster uuid")
		return err
	}

	if id, ok := data[key]; ok {
		initMyCluster(name, id, c_GRU_REMOTE+name)

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

func initMyCluster(name string, id string, remote string) {
	myCluster = Cluster{
		id,
		name,
		remote,
	}
}

func createNodeFolder() error {
	var err error
	opt := discovery.Options{
		"TTL": time.Second * time.Duration(c_TTL),
		"Dir": true,
	}

	config := cfg.GetNodeConfig()
	remote := c_GRU_REMOTE + myCluster.Name + "/" + c_NODES_REMOTE + config.Name
	config.Remote = remote

	err = discovery.Set(remote, "", opt)
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
	remote := cfg.GetNodeConfig().Remote
	err := discovery.Set(remote, "", opt)
	cfg.WriteNodeConfig(remote, cfg.GetNode().Configuration)
	cfg.WriteNodeConstraints(remote, cfg.GetNode().Constraints)
	cfg.WriteNodeResources(remote, cfg.GetNode().Resources)
	cfg.WriteNodeActive(remote, cfg.GetNode().Active)
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
	resp, err := discovery.Get(c_GRU_REMOTE, discovery.Options{})
	if err != nil {
		log.WithField("err", err).Errorln("Error listing clusters")
		return map[string]string{}
	}
	clusters := []string{}
	for k, _ := range resp {
		tokens := strings.Split(k, "/")
		clusters = append(clusters, tokens[2])
	}

	clustersUuid := make(map[string]string, len(clusters))
	for _, name := range clusters {
		resp, err := discovery.Get(c_GRU_REMOTE+name+"/uuid", discovery.Options{})
		if err != nil {
			log.Error("Error getting UUID of cluster ", name)
		}
		clustersUuid[name] = resp[c_GRU_REMOTE+name+"/uuid"]
	}

	return clustersUuid

}

func ListNodes(clusterName string, onlyActive bool) map[string]string {
	nodes := GetNodes(clusterName, onlyActive)
	nameAddress := make(map[string]string, len(nodes))
	for _, node := range nodes {
		nameAddress[node.Configuration.Name] = node.Configuration.Address
	}

	return nameAddress
}

func GetNodes(clusterName string, onlyActive bool) []cfg.Node {
	remote := c_GRU_REMOTE + clusterName + "/nodes"
	nodes := cfg.ReadNodes(remote)

	if onlyActive {
		active := []cfg.Node{}
		for _, node := range nodes {
			if node.Active {
				active = append(active, node)
			}
		}
		log.Debugln("Active nodes: ", active)

		return active
	}

	log.Debugln("Nodes: ", nodes)

	return nodes
}

func GetNode(clusterName string, nodeName string) cfg.Node {
	remote := c_GRU_REMOTE + clusterName + "/nodes/" + nodeName
	return cfg.ReadNode(remote)
}

func ListServices(clusterName string) map[string]string {
	services := GetServices(clusterName)
	list := make(map[string]string, len(services))
	for _, srv := range services {
		list[srv.Name] = srv.Image
	}

	return list
}

func GetServices(clusterName string) []cfg.Service {
	remote := c_GRU_REMOTE + clusterName + "/services"
	return cfg.ReadServices(remote)
}
