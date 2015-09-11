package monitor

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/jbrukh/window"

	"github.com/elleFlorio/gru/node"
)

type GruStats struct {
	Service  map[string]ServiceStats  `json:"service"`
	Instance map[string]InstanceStats `json:"instance"`
	System   SystemStats              `json:"system"`
	Group    GroupStats               `json:"group"`
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

type GroupStats struct {
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

func convertDataToStats(data map[string][]byte) map[string]GruStats {
	stats := make(map[string]GruStats, len(data))
	for key, toConvert := range data {
		converted := GruStats{}
		err := json.Unmarshal(toConvert, &converted)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"key":   key,
			}).Errorln("Error unmarshaling stats data")
		} else {
			stats[key] = converted
		}
	}
	return stats
}

func mergeStats(stats map[string]GruStats) GruStats {
	merged := GruStats{}
	srv_merged := make(map[string]ServiceStats)
	inst_merged := make(map[string]InstanceStats)
	for _, toMerge := range stats {

		//Merge Services and Group stats (for efficency)
		for srv_name, srv_stats := range toMerge.Service {
			srv_stats_merged := srv_merged[srv_name]

			srv_stats_merged.Instances.All = append(srv_stats_merged.Instances.All, srv_stats.Instances.All...)
			merged.Group.Instances.All = append(merged.Group.Instances.All, srv_stats.Instances.All...)

			srv_stats_merged.Instances.Running = append(srv_stats_merged.Instances.Running, srv_stats.Instances.Running...)
			merged.Group.Instances.Running = append(merged.Group.Instances.Running, srv_stats.Instances.Running...)

			srv_stats_merged.Instances.Stopped = append(srv_stats_merged.Instances.Stopped, srv_stats.Instances.Stopped...)
			merged.Group.Instances.Stopped = append(merged.Group.Instances.Stopped, srv_stats.Instances.Stopped...)

			srv_stats_merged.Instances.Paused = append(srv_stats_merged.Instances.Paused, srv_stats.Instances.Paused...)
			merged.Group.Instances.Paused = append(merged.Group.Instances.Paused, srv_stats.Instances.Paused...)

			srv_stats_merged.Events.Start = append(srv_stats_merged.Events.Start, srv_stats.Events.Start...)
			srv_stats_merged.Events.Stop = append(srv_stats_merged.Events.Stop, srv_stats.Events.Stop...)

			srv_merged[srv_name] = srv_stats_merged
		}

		//Merge Instances
		for inst_name, inst_stats := range toMerge.Instance {
			inst_merged[inst_name] = inst_stats
		}
	}

	merged.Service = srv_merged
	merged.Instance = inst_merged
	// System is local status, no need to merge
	merged.System = stats[node.Config().UUID].System

	return merged
}
