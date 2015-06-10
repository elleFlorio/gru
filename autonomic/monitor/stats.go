package monitor

type GruStats struct {
	Service  map[string]ServiceStats
	Instance map[string]InstanceStats
	System   SystemStats
}

type ServiceStats struct {
	Instances []string
	Events    EventStats
}

type EventStats struct {
	Die []string
}

type InstanceStats struct {
	Cpu uint64
}

type SystemStats struct {
	Cpu uint64
}
