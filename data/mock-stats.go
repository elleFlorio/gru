package data

// import (
// 	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/jbrukh/window"
// )

// const maxHistory = 6
// const c_W_SIZE = 100
// const c_W_MULT = 1000

// func ListMockServices() []string {
// 	return []string{
// 		"service1",
// 		"service2",
// 	}
// }

// func StoreMockStats() {
// 	stats := CreateMockStats()
// 	SaveStats(stats)
// }

// func MaxNumberOfEntryInHistory() int {
// 	return maxHistory
// }

// func CreateMockStats() GruStats {
// 	events1 := EventStats{
// 		Stop:  []string{"instance1_0"},
// 		Start: []string{"instance1_4"},
// 	}

// 	cpu1 := CpuStats{
// 		Avg: 2.5,
// 		Tot: 0.7,
// 	}
// 	mem1 := MemoryStats{
// 		Avg: (1 * 1024 * 1024 * 1024),
// 		Tot: (2 * 1024 * 1024 * 1024),
// 	}
// 	metric1 := MetricStats{
// 		map[string]float64{
// 			"CPU_AVG": 1000.0,
// 		},
// 		map[string]float64{
// 			"RESPONSE_TIME": 1000.0,
// 		},
// 	}
// 	service1 := ServiceStats{
// 		Events:  events1,
// 		Cpu:     cpu1,
// 		Memory:  mem1,
// 		Metrics: metric1,
// 	}

// 	cpu2 := CpuStats{
// 		Avg: 0.2,
// 		Tot: 0.2,
// 	}
// 	mem2 := MemoryStats{
// 		Avg: (1 * 1024 * 1024 * 1024),
// 		Tot: (1 * 1024 * 1024 * 1024),
// 	}
// 	metric2 := MetricStats{
// 		map[string]float64{
// 			"CPU_AVG": 2000.0,
// 		},
// 		map[string]float64{
// 			"RESPONSE_TIME": 2000.0,
// 		},
// 	}
// 	service2 := ServiceStats{
// 		Cpu:     cpu2,
// 		Memory:  mem2,
// 		Metrics: metric2,
// 	}
// 	services := map[string]ServiceStats{
// 		"service1": service1,
// 		"service2": service2,
// 	}

// 	cpuSys := 0.8
// 	cpu1_1 := 0.4
// 	mem1_1 := 0.4
// 	cpu1_2 := 0.2
// 	mem1_2 := 0.2
// 	cpu1_3 := 0.0
// 	mem1_3 := 0.0
// 	cpu1_4 := 0.0
// 	mem1_4 := 0.0
// 	cpu2_1 := 0.2
// 	mem2_1 := 0.2

// 	instStat1_1 := InstanceStats{cpu1_1, mem1_1}
// 	instStat1_2 := InstanceStats{cpu1_2, mem1_2}
// 	instStat1_3 := InstanceStats{cpu1_3, mem1_3}
// 	instStat1_4 := InstanceStats{cpu1_4, mem1_4}
// 	instStat2_1 := InstanceStats{cpu2_1, mem2_1}

// 	instances := map[string]InstanceStats{
// 		"instance1_1": instStat1_1,
// 		"instance1_2": instStat1_2,
// 		"instance1_3": instStat1_3,
// 		"instance1_4": instStat1_4,
// 		"instance2_1": instStat2_1,
// 	}

// 	system := SystemStats{
// 		Cpu: cpuSys,
// 	}

// 	mockStats := GruStats{
// 		Service:  services,
// 		Instance: instances,
// 		System:   system,
// 	}

// 	return mockStats
// }

// func CreateMockHistory() StatsHistory {

// 	cpuSysAll := window.New(c_W_SIZE, c_W_MULT)
// 	cpuSysAll.PushBack(float64(1000000))
// 	cpuSysAll.PushBack(float64(1100000))
// 	cpuSysAll.PushBack(float64(1200000))
// 	cpuSysAll.PushBack(float64(1300000))
// 	cpuSysAll.PushBack(float64(1400000))
// 	cpuSysAll.PushBack(float64(1500000))

// 	cpuTot1_1 := window.New(c_W_SIZE, c_W_MULT)
// 	cpuTot1_1.PushBack(float64(10000))
// 	cpuTot1_1.PushBack(float64(20000))
// 	cpuTot1_1.PushBack(float64(30000))
// 	cpuTot1_1.PushBack(float64(40000))
// 	cpuTot1_1.PushBack(float64(50000))
// 	cpuTot1_1.PushBack(float64(60000))
// 	mem1_1 := window.New(c_W_SIZE, c_W_MULT)
// 	mem1_1.PushBack(100000)
// 	mem1_1.PushBack(200000)
// 	mem1_1.PushBack(300000)
// 	mem1_1.PushBack(400000)
// 	mem1_1.PushBack(500000)
// 	mem1_1.PushBack(600000)

// 	cpuTot1_2 := window.New(c_W_SIZE, c_W_MULT)
// 	cpuTot1_2.PushBack(float64(60000))
// 	cpuTot1_2.PushBack(float64(120000))
// 	cpuTot1_2.PushBack(float64(180000))
// 	cpuTot1_2.PushBack(float64(240000))
// 	cpuTot1_2.PushBack(float64(300000))
// 	cpuTot1_2.PushBack(float64(360000))
// 	mem1_2 := window.New(c_W_SIZE, c_W_MULT)
// 	mem1_2.PushBack(150000)
// 	mem1_2.PushBack(250000)
// 	mem1_2.PushBack(350000)
// 	mem1_2.PushBack(450000)
// 	mem1_2.PushBack(550000)
// 	mem1_2.PushBack(650000)

// 	cpuTot1_3 := window.New(c_W_SIZE, c_W_MULT)
// 	cpuTot1_3.PushBack(float64(50000))
// 	cpuTot1_3.PushBack(float64(52000))
// 	cpuTot1_3.PushBack(float64(72000))
// 	cpuTot1_3.PushBack(float64(75000))
// 	cpuTot1_3.PushBack(float64(80000))
// 	cpuTot1_3.PushBack(float64(85000))
// 	mem1_3 := window.New(c_W_SIZE, c_W_MULT)
// 	mem1_3.PushBack(150000)
// 	mem1_3.PushBack(250000)
// 	mem1_3.PushBack(850000)
// 	mem1_3.PushBack(1000000)
// 	mem1_3.PushBack(750000)
// 	mem1_3.PushBack(1200000)

// 	cpuTot1_4 := window.New(c_W_SIZE, c_W_MULT)
// 	cpuTot1_4.PushBack(float64(70000))
// 	mem1_4 := window.New(c_W_SIZE, c_W_MULT)
// 	mem1_4.PushBack(400000)

// 	cpuTot2_1 := window.New(c_W_SIZE, c_W_MULT)
// 	cpuTot2_1.PushBack(float64(40000))
// 	cpuTot2_1.PushBack(float64(80000))
// 	cpuTot2_1.PushBack(float64(120000))
// 	cpuTot2_1.PushBack(float64(160000))
// 	cpuTot2_1.PushBack(float64(200000))
// 	cpuTot2_1.PushBack(float64(240000))
// 	mem2_1 := window.New(c_W_SIZE, c_W_MULT)
// 	mem2_1.PushBack(1000000)
// 	mem2_1.PushBack(1250000)
// 	mem2_1.PushBack(1500000)
// 	mem2_1.PushBack(1750000)
// 	mem2_1.PushBack(2000000)
// 	mem2_1.PushBack(2250000)

// 	cpuHist1_1 := CpuHistory{
// 		cpuTot1_1,
// 		cpuSysAll,
// 	}
// 	cpuHist1_2 := CpuHistory{
// 		cpuTot1_2,
// 		cpuSysAll,
// 	}
// 	cpuHist1_3 := CpuHistory{
// 		cpuTot1_3,
// 		cpuSysAll,
// 	}
// 	cpuHist1_4 := CpuHistory{
// 		cpuTot1_4,
// 		cpuSysAll,
// 	}
// 	cpuHist2_1 := CpuHistory{
// 		cpuTot2_1,
// 		cpuSysAll,
// 	}

// 	instHist1_1 := InstanceHistory{cpuHist1_1, mem1_1}
// 	instHist1_2 := InstanceHistory{cpuHist1_2, mem1_2}
// 	instHist1_3 := InstanceHistory{cpuHist1_3, mem1_3}
// 	instHist1_4 := InstanceHistory{cpuHist1_4, mem1_4}
// 	instHist2_1 := InstanceHistory{cpuHist2_1, mem2_1}

// 	instancesHist := map[string]InstanceHistory{
// 		"instance1_1": instHist1_1,
// 		"instance1_2": instHist1_2,
// 		"instance1_3": instHist1_3,
// 		"instance1_4": instHist1_4,
// 		"instance2_1": instHist2_1,
// 	}

// 	mockHist := StatsHistory{instancesHist}

// 	return mockHist

// }
