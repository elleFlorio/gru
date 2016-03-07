package configuration

type Service struct {
	Name          string             `json:"name"`
	Type          string             `json:"type"`
	Image         string             `json:"image"`
	Remote        string             `json:"remote"`
	DiscoveryPort string             `json:discoveryport`
	Instances     ServiceStatus      `json:"instances"`
	Constraints   ServiceConstraints `json:"constraints"`
	Docker        ServiceDocker      `json:"configuration"`
}

type ServiceStatus struct {
	All     []string `json:"all"`
	Running []string `json:"running"`
	Pending []string `json:"pending"`
	Stopped []string `json:"stopped"`
	Paused  []string `json:"paused"`
}

// TODO this needs a review
type ServiceConstraints struct {
	MaxRespTime float64 `json:"maxresptime"`
}

type ServiceDocker struct {
	Env         map[string]string   `json:"env"`
	Volumes     map[string]struct{} `json:"volumes"`
	Entrypoint  []string            `json:"entrypoint"`
	Memory      string              `json:"memory"`
	CPUnumber   int                 `json:"cpunumber"`
	CpuShares   int64               `json:"cpushares"`
	CpusetCpus  string              `json:"cpusetcpus"`
	Links       []string            `json:"links"`
	Ports       map[string]string   `json:"ports"`
	Cmd         []string            `json:"cmd"`
	StopTimeout int                 `json:"stoptimeout"`
}
