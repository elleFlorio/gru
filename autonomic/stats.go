package autonomic

type GruStats struct {
	Service  map[string]ServiceStats
	Instance map[string]InstanceStats
	System   SystemStats
}

type ServiceStats struct {
	CpuAvg    float64
	Instances []string
}

type InstanceStats struct {
	Cpu uint64
}

type SystemStats struct {
	Cpu uint64
}
