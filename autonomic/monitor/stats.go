package monitor

import (
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/jbrukh/window"

	"github.com/elleFlorio/gru/service"
)

type GruStats struct {
	Service  map[string]ServiceStats  `json:"service"`
	Instance map[string]InstanceStats `json:"instance"`
	System   SystemStats              `json:"system"`
}

type ServiceStats struct {
	Instances service.InstanceStatus `json:"instances"`
	Events    EventStats             `json:"events"`
	Cpu       CpuStats               `json:"cpu"`
	Memory    MemoryStats            `json:memory`
	Metrics   MetricStats            `json:metrics`
}

type EventStats struct {
	Start []string `json:"start"`
	Stop  []string `json:"stop"`
}

type CpuStats struct {
	Avg float64 `json:"avg"`
	Tot float64 `json:"tot"`
}

type MemoryStats struct {
	Avg float64 `json:"avg"`
	Tot float64 `json:"tot"`
}

type MetricStats struct {
	ResponseTime []float64 `json:responsetime`
}

type InstanceStats struct {
	Cpu    float64 `json:"cpu"`
	Memory float64 `json:memory`
}

type SystemStats struct {
	Instances service.InstanceStatus `json:"instances"`
	Cpu       float64                `json:"cpu"`
}

type statsHistory struct {
	instance map[string]instanceHistory
}

type instanceHistory struct {
	cpu cpuHistory
	mem *window.MovingWindow
}

type cpuHistory struct {
	totalUsage *window.MovingWindow
	sysUsage   *window.MovingWindow
}
