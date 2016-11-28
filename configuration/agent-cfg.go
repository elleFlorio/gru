package configuration

type Agent struct {
	Docker        DockerConfig        `json:"docker"`
	Autonomic     AutonomicConfig     `json:"autonomic"`
	Communication CommunicationConfig `json:"communication"`
	Storage       StorageConfig       `json:"storage"`
	Metric        MetricConfig        `json:"metric"`
	Discovery     DiscoveryConfig     `json:"discovery"`
}

type DockerConfig struct {
	DaemonUrl     string `json:"daemonurl"`
	DaemonTimeout int    `json:"daemontimeout"`
}

type AutonomicConfig struct {
	LoopTimeInterval  int    `json:"looptimeinterval"`
	EnableDynamicLoop bool   `json:"enabledynamicloop"`
	PlannerStrategy   string `json:"plannerstrategy"`
	EnableLogReading  bool   `json:"enablelogreading"`
}

type CommunicationConfig struct {
	MaxFriends int `json:"maxfriends"`
}

type StorageConfig struct {
	StorageService string `json:"storageservice"`
}

type MetricConfig struct {
	MetricService    string                 `json:"metricservice"`
	LoopTimeInterval int                    `json:"looptimeinterval"`
	Configuration    map[string]interface{} `json:"configuration"`
}

type DiscoveryConfig struct {
	AppRoot string `json:"approot"`
	TTL     int    `json:ttl`
}
