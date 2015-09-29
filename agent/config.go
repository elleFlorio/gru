package agent

type GruAgentConfig struct {
	Service   ServiceConfig   `json:"service"`
	Node      NodeConfig      `json:"node"`
	Network   NetConfig       `json:"network"`
	Docker    DockerConfig    `json:"docker"`
	Autonomic AutonomicConfig `json:"autonomic"`
	Discovery DiscoveryConfig `json:"discovery"`
	Storage   StorageConfig   `json:"storage"`
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
	MaxFrineds       int    `json:"maxfriends"`
	DataToShare      string `json:"datatoshare"`
}

type DiscoveryConfig struct {
	DiscoveryService    string `json:"discoveryservice"`
	DiscoveryServiceUri string `json:"discoveryurl"`
}

type StorageConfig struct {
	StorageService string `json:"storageservice"`
}
