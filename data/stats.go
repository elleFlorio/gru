package data

import (
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/jbrukh/window"
)

type GruStats struct {
	Metrics MetricStats
	Events  EventStats
}

type MetricStats struct {
	Service  map[string]MetricData `json:"service"`
	Instance map[string]MetricData `json:"instance"`
	System   MetricData            `json:"system"`
}

type MetricData struct {
	BaseMetrics map[string]float64 `json:"basemetrics"`
	UserMetrics map[string]float64 `json:"usermetrics"`
}

type EventStats struct {
	Service map[string]EventData `json:"service"`
}

type EventData struct {
	Start []string `json:"start"`
	Stop  []string `json:"stop"`
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
