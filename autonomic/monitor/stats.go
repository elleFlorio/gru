package monitor

import (
	"encoding/json"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/jbrukh/window"
)

type GruStats struct {
	Service  map[string]ServiceStats  `json:"service"`
	Instance map[string]InstanceStats `json:"instance"`
	System   SystemStats              `json:"system"`
}

type ServiceStats struct {
	Instances InstanceStatus `json:"instances"`
	Events    EventStats     `json:"events"`
	Cpu       CpuStats       `json:"cpu"`
}

type InstanceStatus struct {
	All     []string `json:"all"`
	Running []string `json:"running"`
	Pending []string `json:"pending"`
	Stopped []string `json:"stopped"`
	Paused  []string `json:"paused"`
}

type EventStats struct {
	Start []string `json:"start"`
	Stop  []string `json:"stop"`
}

type CpuStats struct {
	Avg float64 `json:"avg"`
	Tot float64 `json:"tot"`
}

type InstanceStats struct {
	Cpu float64 `json:"cpu"`
}

type SystemStats struct {
	Instances InstanceStatus `json:"instances"`
	Cpu       float64        `json:"cpu"`
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

func convertStatsToData(stats GruStats) ([]byte, error) {
	data, err := json.Marshal(stats)
	if err != nil {
		log.WithField("error", err).Errorln("Error marshaling stats data")
		return nil, err
	}

	return data, nil
}

func ConvertDataToStats(data []byte) (GruStats, error) {
	stats := GruStats{}
	err := json.Unmarshal(data, &stats)
	if err != nil {
		log.WithField("error", err).Errorln("Error unmarshaling stats data")
	}

	return stats, err
}
