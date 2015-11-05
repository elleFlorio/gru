package monitor

import (
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/jbrukh/window"

	"github.com/elleFlorio/gru/service"
)

const maxHistory = 6

func ListMockServices() []string {
	return []string{
		"service1",
		"service2",
	}
}

func StoreMockStats() {
	stats := CreateMockStats()
	saveStats(stats)
}

func MaxNumberOfEntryInHistory() int {
	return maxHistory
}

func CreateMockStats() GruStats {
	all1 := []string{"instance1_0, instance1_1", "instance1_2", "instance1_3", "instance1_4"}
	running1 := []string{"instance1_1", "instance1_2"}
	pending1 := []string{"instance1_3", "instance1_4"}
	stopped1 := []string{"instance1_0"}
	paused1 := []string{}
	events1 := EventStats{
		Stop:  []string{"instance1_0"},
		Start: []string{"instance1_4"},
	}
	instances1 := service.InstanceStatus{
		All:     all1,
		Running: running1,
		Pending: pending1,
		Stopped: stopped1,
		Paused:  paused1,
	}
	cpu1 := CpuStats{
		Avg: 2.5,
		Tot: 0.7,
	}
	service1 := ServiceStats{
		Instances: instances1,
		Events:    events1,
		Cpu:       cpu1,
	}

	all2 := []string{"instance2_1"}
	running2 := []string{"instance2_1"}
	instances2 := service.InstanceStatus{
		All:     all2,
		Running: running2,
	}
	cpu2 := CpuStats{
		Avg: 0.2,
		Tot: 0.2,
	}
	service2 := ServiceStats{Instances: instances2, Cpu: cpu2}
	services := map[string]ServiceStats{
		"service1": service1,
		"service2": service2,
	}

	cpuSys := 0.8
	cpu1_1 := 0.4
	cpu1_2 := 0.2
	cpu1_3 := 0.0
	cpu1_4 := 0.0
	cpu2_1 := 0.2

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
	runningSys := []string{"instance1_1", "instance1_2", "instance2_1"}
	pendingSys := []string{"instance1_3", "instance1_4"}
	stoppedSys := []string{"instance1_0"}
	pausedSys := []string{}
	instancesSys := service.InstanceStatus{
		All:     allSys,
		Running: runningSys,
		Pending: pendingSys,
		Stopped: stoppedSys,
		Paused:  pausedSys,
	}

	system := SystemStats{
		Instances: instancesSys,
		Cpu:       cpuSys,
	}

	mockStats := GruStats{
		Service:  services,
		Instance: instances,
		System:   system,
	}

	return mockStats
}

func CreateMockHistory() statsHistory {

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
	cpuTot1_3.PushBack(float64(80000))
	cpuTot1_3.PushBack(float64(85000))

	cpuTot1_4 := window.New(W_SIZE, W_MULT)
	cpuTot1_4.PushBack(float64(70000))

	cpuTot2_1 := window.New(W_SIZE, W_MULT)
	cpuTot2_1.PushBack(float64(40000))
	cpuTot2_1.PushBack(float64(80000))
	cpuTot2_1.PushBack(float64(120000))
	cpuTot2_1.PushBack(float64(160000))
	cpuTot2_1.PushBack(float64(200000))
	cpuTot2_1.PushBack(float64(240000))

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

	serviceHist := make(map[string]metricsHistory)

	mockHist := statsHistory{serviceHist, instancesHist}

	return mockHist

}
