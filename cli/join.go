package cli

import (
	"encoding/json"
	"fmt"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/codegangsta/cli"

	"github.com/elleFlorio/gru/agent"
	"github.com/elleFlorio/gru/api"
	"github.com/elleFlorio/gru/cluster"
	"github.com/elleFlorio/gru/container"
	"github.com/elleFlorio/gru/discovery"
	"github.com/elleFlorio/gru/metric"
	"github.com/elleFlorio/gru/network"
	"github.com/elleFlorio/gru/node"
	"github.com/elleFlorio/gru/service"
	"github.com/elleFlorio/gru/storage"
	"github.com/elleFlorio/gru/utils"
)

const c_GRU_PATH = "/gru/"
const c_CONFIG_PATH = "config"
const c_CONFIG_SERVICES = "services/"

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

	// infrastructure
	initializeNetwork(ipAddress, port)
	initializeDiscovery("etcd", etcdAddress)
	// Configuration
	initializeAgent(clusterName)
	initializeServices(clusterName)
	// Core agent services
	initializeStorage()
	initializeMetricSerivice()
	initializeContainerEngine()
	// Resources
	initializeNode(nodeName, clusterName)
	// Join Cluster
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

func initializeAgent(clusterName string) {
	configPath := c_GRU_PATH + clusterName + "/" + c_CONFIG_PATH
	agentConfig := agent.GruAgentConfig{}
	getConfig(configPath, &agentConfig)
	agent.Initialize(agentConfig)
}
func initializeServices(clusterName string) {
	services := []service.Service{}
	list := cluster.ListServices(clusterName)
	for _, name := range list {
		servicePath := c_GRU_PATH + clusterName + "/" + c_CONFIG_SERVICES + name
		serviceConfig := service.Service{}
		getConfig(servicePath, &serviceConfig)
		services = append(services, serviceConfig)
	}

	service.Initialize(services)
}

func getConfig(configPath string, config interface{}) {
	var err error

	resp, err := discovery.Get(configPath, discovery.Options{})
	if err != nil {
		log.WithField("err", err).Fatalln("Error getting configuration")
	}

	conf_str := resp[configPath]
	err = json.Unmarshal([]byte(conf_str), config)
	if err != nil {
		log.WithField("err", err).Fatalln("Error unmarshaling configuration")
	}
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

func initializeStorage() {
	_, err := storage.New(agent.Config().Storage.StorageService)
	if err != nil {
		log.WithFields(log.Fields{
			"status":  "warning",
			"error":   err,
			"default": storage.Name(),
		}).Warnln("Error initializing storage service")
	} else {
		log.WithField(storage.Name(), "ok").Infoln("Storage service initialized")
	}
}

func initializeMetricSerivice() {
	_, err := metric.New(agent.Config().Metric.MetricService, agent.Config().Metric.Configuration)
	if err != nil {
		log.WithFields(log.Fields{
			"status":  "warning",
			"error":   err,
			"default": metric.Name(),
		}).Warnln("Error initializing metric service")
	} else {
		log.WithField(metric.Name(), "ok").Infoln("Metric service initialized")
	}
}

func initializeContainerEngine() {
	err := container.Connect(agent.Config().Docker.DaemonUrl, agent.Config().Docker.DaemonTimeout)
	if err != nil {
		log.WithField("err", err).Fatalln("Error initializing container engine")
	}
	log.WithField("docker", "ok").Infoln("Container engine initialized")
}

func initializeNode(nodeName string, clusterName string) {
	if nodeName == "random_name" {
		nodeName = utils.GetRandomName(0)
	}
	counter := -2
	for nameExist(nodeName, clusterName) {
		nodeName = utils.GetRandomName(counter)
		counter++
	}
	log.Debugln("Node name: ", nodeName)
	node.CreateNode(nodeName)
}

func registerToCluster(name string) {
	err := cluster.JoinCluster(name)
	if err != nil {
		log.WithField("err", err).Fatalln("Error registering to cluster")
	}
}
