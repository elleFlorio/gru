package monitor

import (
	"github.com/jbrukh/window"
)

type GruStats struct {
	Service  map[string]ServiceStats  `json:"service"`
	Instance map[string]InstanceStats `json:"instance"`
	System   SystemStats              `json:"system"`
}

type ServiceStats struct {
	Instances InstanceStatus `json:"instances"`
	Events    EventStats     `json:"events"`
}

type InstanceStatus struct {
	All     []string `json:"all"`
	Running []string `json:"running"`
	Stopped []string `json:"stopped"`
	Paused  []string `json:"paused"`
}

type EventStats struct {
	Start []string `json:"start"`
	Stop  []string `json:"stop"`
}

type InstanceStats struct {
	Cpu CpuStats `json:"cpu"`
}

type CpuStats struct {
	TotalUsage []float64 `json:"totalusage"`
	SysUsage   []float64 `json:"sysusage"`
}

type SystemStats struct {
	Instances InstanceStatus `json:"instances"`
}

type statsHistory struct {
	instance map[string]instanceHistory
}

type instanceHistory struct {
	cpu cpuHistory
}

type cpuHistory struct {
	totalUsage *window.MovingWindow
	sysUsage   *window.MovingWindow
}
