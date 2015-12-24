package configuration

type Agent struct {
	Docker        DockerConfig        `json:"docker"`
	Autonomic     AutonomicConfig     `json:"autonomic"`
	Communication CommunicationConfig `json:"communication"`
	Storage       StorageConfig       `json:"storage"`
	Metric        MetricConfig        `json:"metric"`
}

type DockerConfig struct {
	DaemonUrl     string `json:"daemonurl"`
	DaemonTimeout int    `json:"daemontimeout"`
}

type AutonomicConfig struct {
	LoopTimeInterval int    `json:"looptimeinterval"`
	PlannerStrategy  string `json:"plannerstrategy"`
}

type CommunicationConfig struct {
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
