package analyzer

type GruAnalytics struct {
	Service  map[string]ServiceAnalytics  `json:"service"`
	Instance map[string]InstanceAnalytics `json:"instance"`
	System   SystemAnalytics              `json:"system"`
}

type ServiceAnalytics struct {
	CpuTot    float64        `json:"cputot"`
	CpuAvg    float64        `json:"cpuavg"`
	Instances InstanceStatus `json:"instances"`
}

type InstanceStatus struct {
	All     []string `json:"all"`
	Pending []string `json:"pending"`
	Active  []string `json:"active"`
	Stopped []string `json:"stopped"`
	Paused  []string `json:"paused"`
}

type InstanceAnalytics struct {
	Cpu CpuAnalytics `json:"cpu"`
}

type SystemAnalytics struct {
	Cpu       CpuAnalytics   `json:"cpu"`
	Instances InstanceStatus `json:"instances"`
}

type CpuAnalytics struct {
	CpuPerc float64 `json:"cpuperc"`
}
