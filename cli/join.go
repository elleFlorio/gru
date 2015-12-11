package cli

import (
	"fmt"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/codegangsta/cli"

	"github.com/elleFlorio/gru/api"
	"github.com/elleFlorio/gru/cluster"
	"github.com/elleFlorio/gru/discovery"
	"github.com/elleFlorio/gru/network"
	"github.com/elleFlorio/gru/node"
	"github.com/elleFlorio/gru/utils"
)

func join(c *cli.Context) {
	var clusterName string

	if !c.Args().Present() {
		log.Fatalln("Error: missing cluster name")
	} else {
		clusterName = c.Args().First()
	}
	ipAddress := c.String("address")
	port := c.String("port")
	etcdAddress := c.String("etcdserver")
	nodeName := c.String("name")

	initializeNetwork(ipAddress, port)
	initializeDiscovery("etcd", etcdAddress)

	initializeNode(nodeName, clusterName)
	registerToCluster(clusterName)
	defer api.StartServer(port)

	fmt.Printf("Joined cluster %s.\nWaiting for commands...\n", clusterName)
}

func initializeNetwork(address string, port string) {
	err := network.InitializeNetwork(address, port)
	if err != nil {
		log.WithField("err", err).Fatalln("Error initializing the network")
	}
}

func initializeDiscovery(name string, address string) {
	_, err := discovery.New(name, address)
	if err != nil {
		log.WithField("err", err).Fatalln("Error initializing discovery service")
	}
}

func initializeNode(nodeName string, clusterName string) {
	//TODO check for random name
	if nodeName == "random_name" {
		nodeName = utils.GetRandomName(0)
	}
	counter := -2
	for nameExist(nodeName, clusterName) {
		nodeName = utils.GetRandomName(counter)
		counter++
	}

	node.CreateNode(nodeName)
}

func nameExist(nodeName string, clusterName string) bool {
	names := cluster.ListNodes(clusterName, false)
	log.Debugln("Nodes list: ", names)
	for _, name := range names {
		if name == nodeName {
			return true
		}
	}
	return false
}

func registerToCluster(name string) {
	err := cluster.JoinCluster(name)
	if err != nil {
		log.WithField("err", err).Fatalln("Error registering to cluster")
	}
}
