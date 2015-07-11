package analyzer

type GruAnalytics struct {
	Service  map[string]ServiceAnalytics
	Instance map[string]InstanceAnalytics
	System   SystemAnalytics
}

type ServiceAnalytics struct {
	CpuTot    float64
	CpuAvg    float64
	Instances InstanceStatus
}

type InstanceStatus struct {
	All     []string
	Pending []string
	Active  []string
	Stopped []string
	Paused  []string
}

type InstanceAnalytics struct {
	Cpu CpuAnalytics
}

type SystemAnalytics struct {
	Cpu       CpuAnalytics
	Instances InstanceStatus
}

type CpuAnalytics struct {
	CpuPerc float64
}
