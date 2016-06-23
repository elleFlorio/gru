package data

import (
	"github.com/elleFlorio/gru/enum"
)

func CreateMockStats() GruStats {
	// service 1
	srvBaseMetrics1 := map[string]float64{
		enum.METRIC_CPU_AVG.ToString(): 0.6,
		enum.METRIC_MEM_AVG.ToString(): 0.3,
	}
	srvUserMetrics1 := map[string]float64{
		"EXECUTION_TIME": 3000,
	}

	srvMetrics1 := MetricData{
		BaseMetrics: srvBaseMetrics1,
		UserMetrics: srvUserMetrics1,
	}

	// service 2
	srvBaseMetrics2 := map[string]float64{
		enum.METRIC_CPU_AVG.ToString(): 0.1,
		enum.METRIC_MEM_AVG.ToString(): 0.1,
	}
	srvUserMetrics2 := map[string]float64{
		"EXECUTION_TIME": 5000,
	}

	srvMetrics2 := MetricData{
		BaseMetrics: srvBaseMetrics2,
		UserMetrics: srvUserMetrics2,
	}

	// service 3
	srvBaseMetrics3 := map[string]float64{
		enum.METRIC_CPU_AVG.ToString(): 0.9,
		enum.METRIC_MEM_AVG.ToString(): 0.8,
	}
	srvUserMetrics3 := map[string]float64{
		"EXECUTION_TIME": 10000,
	}

	srvMetrics3 := MetricData{
		BaseMetrics: srvBaseMetrics3,
		UserMetrics: srvUserMetrics3,
	}

	srvStats := map[string]MetricData{
		"service1": srvMetrics1,
		"service2": srvMetrics2,
		"service3": srvMetrics3,
	}

	// instance 1
	instBaseMetrics1 := map[string]float64{
		enum.METRIC_CPU_AVG.ToString(): 0.6,
		enum.METRIC_MEM_AVG.ToString(): 0.3,
	}

	instMetrics1 := MetricData{
		BaseMetrics: instBaseMetrics1,
	}

	// instance 2
	instBaseMetrics2 := map[string]float64{
		enum.METRIC_CPU_AVG.ToString(): 0.1,
		enum.METRIC_MEM_AVG.ToString(): 0.1,
	}

	instMetrics2 := MetricData{
		BaseMetrics: instBaseMetrics2,
	}

	// instance 3
	instBaseMetrics3 := map[string]float64{
		enum.METRIC_CPU_AVG.ToString(): 0.9,
		enum.METRIC_MEM_AVG.ToString(): 0.8,
	}

	instMetrics3 := MetricData{
		BaseMetrics: instBaseMetrics3,
	}

	instStats := map[string]MetricData{
		"instance1": instMetrics1,
		"instance2": instMetrics2,
		"instance3": instMetrics3,
	}

	// system
	sysBaseMetrics := map[string]float64{
		enum.METRIC_CPU_AVG.ToString(): 0.5,
		enum.METRIC_MEM_AVG.ToString(): 0.4,
	}

	sysStats := MetricData{
		BaseMetrics: sysBaseMetrics,
	}

	metrics := MetricStats{
		Service:  srvStats,
		Instance: instStats,
		System:   sysStats,
	}

	// events
	events := EventStats{
		Service: map[string]EventData{
			"service1": EventData{},
			"service2": EventData{},
			"service3": EventData{},
		},
	}

	stats := GruStats{
		Metrics: metrics,
		Events:  events,
	}

	return stats
}

func SaveMockStats() {
	stats := CreateMockStats()
	SaveStats(stats)
}
