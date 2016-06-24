package monitor

import "github.com/elleFlorio/gru/data"

func GetStats() data.GruStats {
	return stats
}

func GetServiceStats(name string) data.MetricData {
	return stats.Metrics.Service[name]
}

func GetServicesStats() map[string]data.MetricData {
	return stats.Metrics.Service
}

func GetInstanceStats(id string) data.MetricData {
	return stats.Metrics.Instance[id]
}

func GetInstancesStats() map[string]data.MetricData {
	return stats.Metrics.Instance
}

func GetSystemStats() data.MetricData {
	return stats.Metrics.System
}
