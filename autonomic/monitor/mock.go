package monitor

import "github.com/jbrukh/window"

const maxHistory = 6

func ListMockServices() []string {
	return []string{
		"service1",
		"service2",
	}
}

func MaxNumberOfEntryInHistory() int {
	return maxHistory
}

func CreateMockStats() GruStats {
	all1 := []string{"instance1_0, instance1_1", "instance1_2", "instance1_3", "instance1_4"}
	running1 := []string{"instance1_1", "instance1_2", "instance1_4"}
	stopped1 := []string{"instance1_0"}
	paused1 := []string{"instance1_3"}
	events1 := EventStats{
		Stop:  []string{"instance1_0"},
		Start: []string{"instance1_4"},
	}
	instances1 := InstanceStatus{
		All:     all1,
		Running: running1,
		Stopped: stopped1,
		Paused:  paused1,
	}
	service1 := ServiceStats{
		Instances: instances1,
		Events:    events1,
	}

	all2 := []string{"instance2_1"}
	running2 := []string{"instance2_1"}
	instances2 := InstanceStatus{
		All:     all2,
		Running: running2,
	}
	service2 := ServiceStats{Instances: instances2}
	services := map[string]ServiceStats{
		"service1": service1,
		"service2": service2,
	}

	cpuSysAll := []float64{1000000, 1100000, 1200000, 1300000, 1400000, 1500000}
	cpuTot1_1 := []float64{10000, 20000, 30000, 40000, 50000, 60000}
	cpuTot1_2 := []float64{60000, 120000, 180000, 240000, 300000, 360000}
	cpuTot1_3 := []float64{50000, 52000, 72000, 75000}
	cpuTot1_4 := []float64{70000}
	cpuTot2_1 := []float64{40000, 80000, 120000, 160000, 200000, 240000}

	cpu1_1 := CpuStats{
		cpuTot1_1,
		cpuSysAll,
	}
	cpu1_2 := CpuStats{
		cpuTot1_2,
		cpuSysAll,
	}
	cpu1_3 := CpuStats{
		cpuTot1_3,
		cpuSysAll[0:4],
	}
	cpu1_4 := CpuStats{
		cpuTot1_4,
		cpuSysAll[:1],
	}
	cpu2_1 := CpuStats{
		cpuTot2_1,
		cpuSysAll,
	}

	instStat1_1 := InstanceStats{cpu1_1}
	instStat1_2 := InstanceStats{cpu1_2}
	instStat1_3 := InstanceStats{cpu1_3}
	instStat1_4 := InstanceStats{cpu1_4}
	instStat2_1 := InstanceStats{cpu2_1}

	instances := map[string]InstanceStats{
		"instance1_1": instStat1_1,
		"instance1_2": instStat1_2,
		"instance1_3": instStat1_3,
		"instance1_4": instStat1_4,
		"instance2_1": instStat2_1,
	}

	allSys := []string{"instance1_0, instance1_1", "instance1_2", "instance1_3", "instance1_4", "instance2_1"}
	runningSys := []string{"instance1_1", "instance1_2", "instance1_4", "instance2_1"}
	stoppedSys := []string{"instance1_0"}
	pausedSys := []string{"instance1_3"}
	instancesSys := InstanceStatus{
		All:     allSys,
		Running: runningSys,
		Stopped: stoppedSys,
		Paused:  pausedSys,
	}

	system := SystemStats{
		Instances: instancesSys,
	}

	mockStats := GruStats{
		Service:  services,
		Instance: instances,
		System:   system,
	}

	return mockStats
}

func CreateMockHistory() *statsHistory {

	cpuSysAll := window.New(W_SIZE, W_MULT)
	cpuSysAll.PushBack(float64(1000000))
	cpuSysAll.PushBack(float64(1100000))
	cpuSysAll.PushBack(float64(1200000))
	cpuSysAll.PushBack(float64(1300000))
	cpuSysAll.PushBack(float64(1400000))
	cpuSysAll.PushBack(float64(1500000))

	cpuTot1_1 := window.New(W_SIZE, W_MULT)
	cpuTot1_1.PushBack(float64(10000))
	cpuTot1_1.PushBack(float64(20000))
	cpuTot1_1.PushBack(float64(30000))
	cpuTot1_1.PushBack(float64(40000))
	cpuTot1_1.PushBack(float64(50000))
	cpuTot1_1.PushBack(float64(60000))

	cpuTot1_2 := window.New(W_SIZE, W_MULT)
	cpuTot1_2.PushBack(float64(60000))
	cpuTot1_2.PushBack(float64(120000))
	cpuTot1_2.PushBack(float64(180000))
	cpuTot1_2.PushBack(float64(240000))
	cpuTot1_2.PushBack(float64(300000))
	cpuTot1_2.PushBack(float64(360000))

	cpuTot1_3 := window.New(W_SIZE, W_MULT)
	cpuTot1_3.PushBack(float64(50000))
	cpuTot1_3.PushBack(float64(52000))
	cpuTot1_3.PushBack(float64(72000))
	cpuTot1_3.PushBack(float64(75000))

	cpuTot1_4 := window.New(W_SIZE, W_MULT)
	cpuTot1_4.PushBack(float64(70000))

	cpuTot2_1 := window.New(W_SIZE, W_MULT)
	cpuSysAll.PushBack(float64(40000))
	cpuSysAll.PushBack(float64(80000))
	cpuSysAll.PushBack(float64(120000))
	cpuSysAll.PushBack(float64(160000))
	cpuSysAll.PushBack(float64(200000))
	cpuSysAll.PushBack(float64(240000))

	cpuHist1_1 := cpuHistory{
		cpuTot1_1,
		cpuSysAll,
	}
	cpuHist1_2 := cpuHistory{
		cpuTot1_2,
		cpuSysAll,
	}
	cpuHist1_3 := cpuHistory{
		cpuTot1_3,
		cpuSysAll,
	}
	cpuHist1_4 := cpuHistory{
		cpuTot1_4,
		cpuSysAll,
	}
	cpuHist2_1 := cpuHistory{
		cpuTot2_1,
		cpuSysAll,
	}

	instHist1_1 := instanceHistory{cpuHist1_1}
	instHist1_2 := instanceHistory{cpuHist1_2}
	instHist1_3 := instanceHistory{cpuHist1_3}
	instHist1_4 := instanceHistory{cpuHist1_4}
	instHist2_1 := instanceHistory{cpuHist2_1}

	instancesHist := map[string]instanceHistory{
		"instance1_1": instHist1_1,
		"instance1_2": instHist1_2,
		"instance1_3": instHist1_3,
		"instance1_4": instHist1_4,
		"instance2_1": instHist2_1,
	}

	mockHist := statsHistory{instancesHist}

	return &mockHist

}
