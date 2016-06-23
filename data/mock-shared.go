package data

import (
	"github.com/elleFlorio/gru/enum"
)

func CreateMockShared() Shared {
	// service 1
	srvBaseShared1 := map[string]float64{
		enum.METRIC_CPU_AVG.ToString(): 0.6,
		enum.METRIC_MEM_AVG.ToString(): 0.3,
	}
	srvUserShared1 := map[string]float64{
		"LOAD": 0.5,
	}

	srvSharedData1 := SharedData{
		BaseShared: srvBaseShared1,
		UserShared: srvUserShared1,
	}

	srvShared1 := ServiceShared{
		Data:   srvSharedData1,
		Active: true,
	}

	// service 2
	srvBaseShared2 := map[string]float64{
		enum.METRIC_CPU_AVG.ToString(): 0.1,
		enum.METRIC_MEM_AVG.ToString(): 0.1,
	}
	srvUserShared2 := map[string]float64{
		"LOAD": 0.1,
	}

	srvSharedData2 := SharedData{
		BaseShared: srvBaseShared2,
		UserShared: srvUserShared2,
	}

	srvShared2 := ServiceShared{
		Data:   srvSharedData2,
		Active: true,
	}

	// service 3
	srvBaseShared3 := map[string]float64{
		enum.METRIC_CPU_AVG.ToString(): 0.9,
		enum.METRIC_MEM_AVG.ToString(): 0.8,
	}
	srvUserShared3 := map[string]float64{
		"LOAD": 0.9,
	}

	srvSharedData3 := SharedData{
		BaseShared: srvBaseShared3,
		UserShared: srvUserShared3,
	}

	srvShared3 := ServiceShared{
		Data:   srvSharedData3,
		Active: true,
	}

	// system
	sysBaseShared := map[string]float64{
		enum.METRIC_CPU_AVG.ToString(): 0.5,
		enum.METRIC_MEM_AVG.ToString(): 0.4,
	}

	sysSharedData := SharedData{
		BaseShared: sysBaseShared,
	}

	sysShared := SystemShared{
		Data:           sysSharedData,
		ActiveServices: []string{"service1", "service2", "service3"},
	}

	shared := Shared{
		Service: map[string]ServiceShared{
			"service1": srvShared1,
			"service2": srvShared2,
			"service3": srvShared3,
		},
		System: sysShared,
	}

	return shared
}
