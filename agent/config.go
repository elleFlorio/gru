package agent

type GruAgentConfig struct {
	Service   ServiceConfig   `json:"service"`
	Node      NodeConfig      `json:"node"`
	Network   NetConfig       `json:"network"`
	Docker    DockerConfig    `json:"docker"`
	Autonomic AutonomicConfig `json:"autonomic"`
	Discovery DiscoveryConfig `json:"discovery"`
	Storage   StorageConfig   `json:"storage"`
	Metric    MetricConfig    `json:"metric"`
}

type ServiceConfig struct {
	ServiceConfigFolder string `json:"serviceconfigfolder"`
}

type NodeConfig struct {
	NodeConfigFile string `json:"nodeconfigfile"`
}

type NetConfig struct {
	IpAddres string `json:"ipaddress"`
	Port     string `json:"port"`
}

type DockerConfig struct {
	DaemonUrl     string `json:"daemonurl"`
	DaemonTimeout int    `json:"daemontimeout"`
}

type AutonomicConfig struct {
	LoopTimeInterval int    `json:"looptimeinterval"`
	MaxFriends       int    `json:"maxfriends"`
	DataToShare      string `json:"datatoshare"`
}

type DiscoveryConfig struct {
	DiscoveryService    string `json:"discoveryservice"`
	DiscoveryServiceUri string `json:"discoveryserviceuri"`
}

type StorageConfig struct {
	StorageService string `json:"storageservice"`
}

// This is a temporal solution.
// When and if new services will be added I will find
// the correct way to make it generic
type MetricConfig struct {
	MetricService string                 `json:"metricservice"`
	Configuration map[string]interface{} `json:configuration`
}
