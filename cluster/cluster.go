package cluster

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/discovery"
	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/network"
	"github.com/elleFlorio/gru/node"
	"github.com/elleFlorio/gru/storage"
)

// etcd
const c_GRU_PATH = "/gru/"
const c_NODES_FOLDER = "nodes/"
const c_CONFIG_FOLDER = "config/"
const c_SERVICES_FOLDER = "services/"
const c_TTL = 5

// api
const c_ROUTE_ANALYTICS string = "/gru/v1/analytics"

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
		initCluster(name, id, c_GRU_PATH+name, c_GRU_PATH+name+"/"+c_NODES_FOLDER+node.Config().Name)
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
	err = discovery.Set(myCluster.NodePath, "", opt)
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
	err = discovery.Set(myCluster.NodePath, "", opt)
	err = discovery.Set(myCluster.NodePath+"/uuid", node.Config().UUID, discovery.Options{})
	err = discovery.Set(myCluster.NodePath+"/name", node.Config().Name, discovery.Options{})
	err = discovery.Set(myCluster.NodePath+"/address", node.Config().Address, discovery.Options{})
	err = discovery.Set(myCluster.NodePath+"/active", active, discovery.Options{})
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

func ListNodes(clusterName string, onlyActive bool) []string {
	nodesPath := c_GRU_PATH + clusterName + "/" + c_NODES_FOLDER
	resp, err := discovery.Get(nodesPath, discovery.Options{"Recursive": true})
	if err != nil {
		log.WithField("err", err).Errorln("Error listing nodes in cluster ", clusterName)
		return []string{}
	}
	return getNodeNames(resp, onlyActive, nodesPath)
}

func GetNodes(clusterName string, onlyActive bool) []node.Node {
	nodes := []node.Node{}
	nodesPath := c_GRU_PATH + clusterName + "/nodes/"
	resp, err := discovery.Get(nodesPath, discovery.Options{"Recursive": true})
	if err != nil {
		log.WithField("err", err).Errorln("Error listing nodes in cluster ", clusterName)
		return nodes
	}

	names := getNodeNames(resp, onlyActive, nodesPath)
	for _, name := range names {
		newNode := node.Node{}
		newNode.Name = name
		newNode.UUID = resp[nodesPath+name+"/uuid"]
		newNode.Address = resp[nodesPath+name+"/address"]
		newNode.Active, _ = strconv.ParseBool(resp[nodesPath+name+"/active"])
		if onlyActive {
			if newNode.Active {
				nodes = append(nodes, newNode)
			}
		} else {
			nodes = append(nodes, newNode)
		}
	}

	log.Debugln("Nodes: ", nodes)

	return nodes
}

func getNodeNames(resp map[string]string, onlyActive bool, nodesPath string) []string {
	names := []string{}
	for key, _ := range resp {
		path := strings.Split(key, "/")
		fieldName := path[len(path)-1]

		if fieldName == "nodes" {
			log.Debugln("Nodes folder empty")
			return names
		}

		if fieldName == "name" {
			nodeName := resp[key]
			if onlyActive {
				active, _ := strconv.ParseBool(resp[nodesPath+nodeName+"/active"])
				if active {
					names = append(names, nodeName)
				}
			} else {
				names = append(names, nodeName)
			}
		}
	}

	return names
}

func ListServices(clusterName string) []string {
	servicesPath := c_GRU_PATH + clusterName + "/" + c_SERVICES_FOLDER
	resp, err := discovery.Get(servicesPath, discovery.Options{})
	if err != nil {
		log.WithField("err", err).Errorln("Error listing services in cluster ", clusterName)
		return []string{}
	}

	services := []string{}
	for k, _ := range resp {
		path := strings.Split(k, "/")
		name := path[len(path)-1]
		if name == "services" {
			log.Debugln("Services folder empty")
			return services
		}

		services = append(services, name)
	}

	return services
}

func UpdateFriendsData(nFriends int) error {
	log.Debugln("Updating friends data")
	storage.DeleteAllData(enum.ANALYTICS)

	peers := getAllPeers()
	log.WithField("peers", len(peers)).Debugln("Number of peers")
	if len(peers) == 0 {
		return ErrNoPeers
	}

	friends, err := chooseRandomFriends(peers, nFriends)
	if err != nil {
		return err
	}
	log.WithField("friends", friends).Debugln("Friends to connect with")

	err = getFriendsData(friends)
	if err != nil {
		return err
	}

	return nil
}

func getAllPeers() map[string]string {
	peers := map[string]string{}
	nodes := GetNodes(myCluster.Name, true)
	for _, peer := range nodes {
		peers[peer.Name] = peer.Address
	}

	return peers
}

// Is there a more efficient way to do this?
func chooseRandomFriends(peers map[string]string, n int) (map[string]string, error) {
	nPeers := len(peers)
	if nPeers <= 1 {
		return nil, ErrNoFriends
	}

	if n <= 0 {
		return nil, ErrInvalidFriendsNumber
	} else if n > nPeers {
		n = nPeers
	}

	friends := make(map[string]string, n)

	peersKeys := make([]string, 0, len(peers))
	for peerKey, _ := range peers {
		peersKeys = append(peersKeys, peerKey)
	}

	friendsKeys := make([]string, 0, n)
	indexes := rand.Perm(nPeers)[:n]
	for _, index := range indexes {
		if peersKeys[index] != node.Config().UUID {
			friendsKeys = append(friendsKeys, peersKeys[index])
		}
	}

	for _, friendKey := range friendsKeys {
		friends[friendKey] = peers[friendKey]
	}

	return friends, nil
}

func getFriendsData(friends map[string]string) error {
	var err error

	for friend, address := range friends {
		friendRoute := address + c_ROUTE_ANALYTICS
		friendData, err := network.DoRequest("GET", friendRoute, nil)
		if err != nil {
			log.WithField("address", address).Debugln("Error retrieving friend stats")
		}
		err = storage.StoreData(friend, friendData, enum.ANALYTICS)
		if err != nil {
			log.WithField("err", err).Debugln("Error storing friend stats")
		}
	}

	return err
}
