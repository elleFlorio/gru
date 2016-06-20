package monitor

import "github.com/elleFlorio/gru/data"

func GetStats() data.GruStats {
	return stats
}

func GetServiceMetrics(name string) data.MetricData {
	return stats.Metrics.Service[name]
}

func GetServicesMetrics() map[string]data.MetricData {
	return stats.Metrics.Service
}

func GetInstanceMetrics(id string) data.MetricData {
	return stats.Metrics.Instance[id]
}

func GetInstancesStats() map[string]data.MetricData {
	return stats.Metrics.Instance
}

func GetSystemMetrics() data.MetricData {
	return stats.Metrics.System
}
