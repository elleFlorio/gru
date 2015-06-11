package monitor

type GruStats struct {
	Service  map[string]ServiceStats
	Instance map[string]InstanceStats
	System   SystemStats
}

type ServiceStats struct {
	Instances InstanceStatus
	Events    EventStats
}

type InstanceStatus struct {
	All     []string
	Running []string
	Stopped []string
	Paused  []string
}

type EventStats struct {
	Start []string
	Stop  []string
}

type InstanceStats struct {
	Cpu uint64
}

type SystemStats struct {
	Cpu uint64
}
