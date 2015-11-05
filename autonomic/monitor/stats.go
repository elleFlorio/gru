package monitor

import (
	"encoding/json"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
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

type MetricStats struct {
	ResponseTime []float64 `json:responsetime`
}

type InstanceStats struct {
	Cpu float64 `json:"cpu"`
}

type SystemStats struct {
	Instances service.InstanceStatus `json:"instances"`
	Cpu       float64                `json:"cpu"`
}

type statsHistory struct {
	service  map[string]metricsHistory //Deprecated?
	instance map[string]instanceHistory
}

type metricsHistory struct {
	responseTime *window.MovingWindow
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
