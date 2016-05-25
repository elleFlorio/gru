package cli

import (
	"fmt"
	"strings"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/codegangsta/cli"

	"github.com/elleFlorio/gru/agent"
	"github.com/elleFlorio/gru/api"
	"github.com/elleFlorio/gru/cluster"
	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/container"
	"github.com/elleFlorio/gru/discovery"
	"github.com/elleFlorio/gru/metric"
	"github.com/elleFlorio/gru/network"
	"github.com/elleFlorio/gru/node"
	res "github.com/elleFlorio/gru/resources"
	"github.com/elleFlorio/gru/storage"
	"github.com/elleFlorio/gru/utils"
)

const c_GRU_REMOTE = "/gru/"
const c_CONFIG_REMOTE = "config"
const c_SERVICES_REMOTE = "services"
const c_TUNING_REMOTE = "tuning"

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
	initializeTuning(clusterName)
	// Core agent services
	initializeStorage()
	initializeMetricSerivice()
	initializeContainerEngine()
	// Resources
	initializeResources()
	initializeNode(nodeName, clusterName)
	// Join Cluster
	registerToCluster(clusterName)
	agent.StartMonitoring()
	defer api.StartServer(port)

	fmt.Printf("Joined cluster %s.\nWaiting for commands...\n", clusterName)
}

func initializeNetwork(address string, port string) {
	err := network.InitializeNetwork(address, port)
	if err != nil {
		log.WithField("err", err).Fatalln("Error initializing the network")
	}
	log.WithField(network.Config().IpAddress+":"+network.Config().Port, "ok").Infoln("Network initialized")
}

func initializeDiscovery(name string, address string) {
	_, err := discovery.New(name, address)
	if err != nil {
		log.WithField("err", err).Fatalln("Error initializing discovery service")
	}
}

func initializeAgent(clusterName string) {
	configPath := c_GRU_REMOTE + clusterName + "/" + c_CONFIG_REMOTE
	agentConfig := cfg.Agent{}
	cfg.ReadAgentConfig(configPath, &agentConfig)
	cfg.SetAgent(agentConfig)
	agent.Initialize()
}

func initializeServices(clusterName string) {
	remote := c_GRU_REMOTE + clusterName + "/" + c_SERVICES_REMOTE
	services := cfg.ReadServices(remote)
	cfg.SetServices(services)
}

func initializeTuning(clusterName string) {
	remote := c_GRU_REMOTE + clusterName + "/" + c_TUNING_REMOTE
	tuning := cfg.ReadTuningConfig(remote)
	cfg.SetTuning(tuning)
	log.WithField("tuning", tuning).Debugln("Tuning read from remote")
}

func initializeStorage() {
	_, err := storage.New(cfg.GetAgentStorage().StorageService)
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
	_, err := metric.New(cfg.GetAgentMetric().MetricService, cfg.GetAgentMetric().Configuration)
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

// TODO needs to make more "generic" the container engine package
func initializeContainerEngine() {
	daemonUrl := getDaemonUrl()
	log.WithField("daemonUrl", daemonUrl).Debugln("Container engine initialization")
	err := container.Connect(daemonUrl, cfg.GetAgentDocker().DaemonTimeout)
	if err != nil {
		log.WithField("err", err).Fatalln("Error initializing container engine")
	}
	log.WithField("docker", "ok").Infoln("Container engine initialized")
}

func getDaemonUrl() string {
	url := cfg.GetAgentDocker().DaemonUrl
	if url == "default" {
		return "unix:///var/run/docker.sock"
	}

	ipPort := strings.Split(url, ":")
	if len(ipPort) < 2 {
		log.WithField("err", "Wrong daemon url").Fatalln("Error initializing container engine")
	}

	var daemonUrl string
	if ipPort[0] == "local" {
		daemonUrl = "http://" + network.Config().IpAddress + ":" + ipPort[1]
	} else {
		daemonUrl = "http://" + ipPort[0] + ":" + ipPort[1]
	}

	return daemonUrl
}

func initializeResources() {
	res.Initialize()
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
	node.CreateNode(nodeName, res.GetResources())
}

func nameExist(nodeName string, clusterName string) bool {
	names := cluster.ListNodes(clusterName, false)
	log.Debugln("Nodes list: ", names)
	for name, _ := range names {
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
