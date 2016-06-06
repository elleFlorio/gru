package data

import (
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/jbrukh/window"
)

type GruStats struct {
	Service  map[string]ServiceStats  `json:"service"`
	Instance map[string]InstanceStats `json:"instance"`
	System   SystemStats              `json:"system"`
}

type ServiceStats struct {
	Events  EventStats  `json:"events"`
	Cpu     CpuStats    `json:"cpu"`
	Memory  MemoryStats `json:memory`
	Metrics MetricStats `json:metrics`
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
	BaseMetrics map[string]float64 `json:basemetrics`
	UserMetrics map[string]float64 `json:usermetrics`
}

type InstanceStats struct {
	Cpu    float64 `json:"cpu"`
	Memory float64 `json:memory`
}

type SystemStats struct {
	Cpu float64 `json:"cpu"`
}

type StatsHistory struct {
	Instance map[string]InstanceHistory
}

type InstanceHistory struct {
	Cpu CpuHistory
	Mem *window.MovingWindow
}

type CpuHistory struct {
	TotalUsage *window.MovingWindow
	SysUsage   *window.MovingWindow
}
