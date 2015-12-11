package agent

type GruAgentConfig struct {
	Docker    DockerConfig    `json:"docker"`
	Autonomic AutonomicConfig `json:"autonomic"`
	Storage   StorageConfig   `json:"storage"`
	Metric    MetricConfig    `json:"metric"`
}

type DockerConfig struct {
	DaemonUrl     string `json:"daemonurl"`
	DaemonTimeout int    `json:"daemontimeout"`
}

type AutonomicConfig struct {
	LoopTimeInterval int `json:"looptimeinterval"`
	MaxFriends       int `json:"maxfriends"`
}

type StorageConfig struct {
	StorageService string `json:"storageservice"`
}

type MetricConfig struct {
	MetricService string                 `json:"metricservice"`
	Configuration map[string]interface{} `json:"configuration"`
}
