package monitor

import (
	"github.com/jbrukh/window"
)

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
	Cpu CpuStats
}

type CpuStats struct {
	TotalUsage []float64
	SysUsage   []float64
}

type SystemStats struct {
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
