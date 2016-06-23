package data

import (
	"github.com/elleFlorio/gru/enum"
)

func CreateMockAnalytics() GruAnalytics {
	// service 1
	srvBaseAnalytics1 := map[string]float64{
		enum.METRIC_CPU_AVG.ToString(): 0.6,
		enum.METRIC_MEM_AVG.ToString(): 0.3,
	}
	srvUserAnalytics1 := map[string]float64{
		"LOAD": 0.5,
	}

	srvAnalytics1 := AnalyticData{
		BaseAnalytics: srvBaseAnalytics1,
		UserAnalytics: srvUserAnalytics1,
	}

	// service 2
	srvBaseAnalytics2 := map[string]float64{
		enum.METRIC_CPU_AVG.ToString(): 0.1,
		enum.METRIC_MEM_AVG.ToString(): 0.1,
	}
	srvUserAnalytics2 := map[string]float64{
		"LOAD": 0.1,
	}

	srvAnalytics2 := AnalyticData{
		BaseAnalytics: srvBaseAnalytics2,
		UserAnalytics: srvUserAnalytics2,
	}

	// service 3
	srvBaseAnalytics3 := map[string]float64{
		enum.METRIC_CPU_AVG.ToString(): 0.9,
		enum.METRIC_MEM_AVG.ToString(): 0.8,
	}
	srvUserAnalytics3 := map[string]float64{
		"LOAD": 0.9,
	}

	srvAnalytics3 := AnalyticData{
		BaseAnalytics: srvBaseAnalytics3,
		UserAnalytics: srvUserAnalytics3,
	}

	srvAnalytics := map[string]AnalyticData{
		"service1": srvAnalytics1,
		"service2": srvAnalytics2,
		"service3": srvAnalytics3,
	}

	// system
	sysBaseAnalytics := map[string]float64{
		enum.METRIC_CPU_AVG.ToString(): 0.5,
		enum.METRIC_MEM_AVG.ToString(): 0.4,
	}

	sysAnalytics := AnalyticData{
		BaseAnalytics: sysBaseAnalytics,
	}

	analytics := GruAnalytics{
		Service: srvAnalytics,
		System:  sysAnalytics,
	}

	return analytics
}

func SaveMockAnalytics() {
	analytics := CreateMockAnalytics()
	SaveAnalytics(analytics)
}
