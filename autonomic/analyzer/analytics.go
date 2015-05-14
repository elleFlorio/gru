package analyzer

type GruAnalytics struct {
	Service  map[string]ServiceAnalytics
	Instance map[string]InstanceAnalytics
	System   SystemAnalytics
}

type ServiceAnalytics struct {
	CpuAvg    float64
	Instances []string
}

type InstanceAnalytics struct {
	Cpu     uint64
	CpuPerc float64
}

type SystemAnalytics struct {
	Cpu uint64
}
