package cluster

import (
	"errors"
	"strings"
	"time"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/discovery"
	"github.com/elleFlorio/gru/node"
)

const c_GRU_PATH = "/gru/"
const c_TTL = 5

type Cluster struct {
	UUID        string
	Name        string
	ClusterPath string
	NodesPath   string
}

var (
	myCluster Cluster

	ErrNoClusterId = errors.New("Cannot find cluster ID")
)

func GetMyCluster() Cluster {
	return myCluster
}

func RegisterCluster(name string, id string) {
	err := discovery.Register(c_GRU_PATH+name+"/uuid", id)
	if err != nil {
		log.Errorln("Error registering cluster")
	}
}

func JoinCluster(name string) error {
	key := c_GRU_PATH + name + "/uuid"
	data, err := discovery.Get(key, discovery.Options{})
	if err != nil {
		log.WithField("err", err).Errorln("Error getting cluster uuid")
		return err
	}
	if id, ok := data[key]; ok {
		initCluster(name, id, c_GRU_PATH+name, c_GRU_PATH+name+"/nodes/"+node.Config().Name)
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

func initCluster(name string, id string, clusterPath string, nodesPath string) {
	myCluster = Cluster{
		id,
		name,
		clusterPath,
		nodesPath,
	}
}

func keepAlive(ttl int) {
	ticker := time.NewTicker(time.Second * time.Duration(ttl))
	for {
		select {
		case <-ticker.C:
			err := updateNodeFolder(ttl)
			log.Debugln("Agent is alive")
			if err != nil {
				log.Errorln("Error keeping the node alive")
			}
		}
	}
}

func createNodeFolder() error {
	var err error
	opt := discovery.Options{
		"TTL": time.Second * time.Duration(c_TTL),
		"Dir": true,
	}
	err = discovery.Set(myCluster.NodesPath, "", opt)
	if err != nil {
		log.WithField("err", err).Errorln("Error creating node folder")
		return err
	}

	return nil
}

func updateNodeFolder(ttl int) error {
	var err error
	var active string
	if node.Config().Active {
		active = "true"
	} else {
		active = "false"
	}

	opt := discovery.Options{
		"TTL":       time.Second * time.Duration(ttl),
		"Dir":       true,
		"PrevExist": true,
	}
	err = discovery.Set(myCluster.NodesPath, "", opt)
	err = discovery.Set(myCluster.NodesPath+"/uuid", node.Config().UUID, discovery.Options{})
	err = discovery.Set(myCluster.NodesPath+"/name", node.Config().Name, discovery.Options{})
	err = discovery.Set(myCluster.NodesPath+"/address", node.Config().Address, discovery.Options{})
	err = discovery.Set(myCluster.NodesPath+"/active", active, discovery.Options{})
	if err != nil {
		log.WithField("err", err).Errorln("Error updating node folder")
		return err
	}

	return nil
}

func ListNodes(clusterName string, onlyActive bool) []string {
	nodes := []string{}
	nodesPath := c_GRU_PATH + clusterName + "/nodes"
	data, err := discovery.Get(nodesPath, discovery.Options{})
	if err != nil {
		log.WithField("err", err).Errorln("Error listing nodes in cluster ", clusterName)
		return nodes
	}

	for key, _ := range data {
		path := strings.Split(key, "/")
		name := path[len(path)-1]
		if onlyActive {
			if isActiveNode(clusterName, name) {
				nodes = append(nodes, name)
			}
		} else {
			nodes = append(nodes, name)
		}
	}

	log.Debugln("Node List: ", nodes)

	return nodes
}

//TODO this is very inefficient... Is there a better way?
func isActiveNode(clusterName string, nodeName string) bool {
	activePath := c_GRU_PATH + clusterName + "/nodes/" + nodeName + "/active"
	data, err := discovery.Get(activePath, discovery.Options{})
	if err != nil {
		log.WithField("err", err).Errorln("Error checking active node in cluster ", clusterName)
		return false
	}

	if data[activePath] == "true" {
		return true
	}

	return false
}
