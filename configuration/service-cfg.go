package configuration

type Service struct {
	Name        string             `json:"name"`
	Type        string             `json:"type"`
	Image       string             `json:"image"`
	Remote      string             `json:"remote"`
	Instances   ServiceStatus      `json:"instances"`
	Constraints ServiceConstraints `json:"constraints"`
	Docker      ServiceDocker      `json:"configuration"`
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
	Cmd          []string                 `json:"cmd"`
	Volumes      map[string]struct{}      `json:"volumes"`
	Entrypoint   []string                 `json:"entrypoint"`
	Memory       string                   `json:"memory"`
	CpuShares    int64                    `json:"cpushares"`
	CpusetCpus   string                   `json:"cpusetcpus"`
	ExposedPorts map[string]struct{}      `json:"exposedports"`
	PortBindings map[string][]PortBinding `json:"portbindings"`
	Links        []string                 `json:"links"`
	StopTimeout  int                      `json:"stoptimeout"`
}

type PortBinding struct {
	HostIp   string `json:"hostip"`
	HostPort string `json:"hostport"`
}
