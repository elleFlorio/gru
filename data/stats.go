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
	Metrics MetricStats `json:metrics`
}

type EventStats struct {
	Start []string `json:"start"`
	Stop  []string `json:"stop"`
}

type MetricStats struct {
	BaseMetrics map[string]float64 `json:basemetrics`
	UserMetrics map[string]float64 `json:usermetrics`
}

type InstanceStats struct {
	BaseMetrics map[string]float64 `json:basemetrics`
}

type SystemStats struct {
	BaseMetrics map[string]float64 `json:basemetrics`
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
