package service

type Service struct {
	Name          string         `json:"name"`
	Type          string         `json:"type"`
	Image         string         `json:"image"`
	Instances     InstanceStatus `json:"instances"`
	Constraints   Constraints    `json:"constraints"` //Needed?
	Configuration Config         `json: "configuration"`
}

type InstanceStatus struct {
	All     []string `json:"all"`
	Running []string `json:"running"`
	Pending []string `json:"pending"`
	Stopped []string `json:"stopped"`
	Paused  []string `json:"paused"`
}

// TODO this needs a review
type Constraints struct {
	MaxRespTime float64 `json:"maxresptime"`
}

type Config struct {
	Cmd          []string                 `json:"cmd"`
	Volumes      map[string]struct{}      `json:"volumes"`
	Entrypoint   []string                 `json:"entrypoint"`
	Memory       string                   `json:"memory"`
	CpuShares    int64                    `json:"cpushares"`
	CpusetCpus   string                   `json:"cpusetcpus"`
	PortBindings map[string][]PortBinding `json:"portbindings"`
	Links        []string                 `json:"links"`
}

type PortBinding struct {
	HostIp   string `json:"hostip"`
	HostPort string `json:"hostport"`
}
