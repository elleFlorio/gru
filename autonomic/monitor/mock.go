package monitor

import "github.com/jbrukh/window"

func ListMockServices() []string {
	return []string{
		"service1",
		"service2",
	}
}

func CreateMockStats() GruStats {
	all1 := []string{"instance1_1", "instance1_2", "instance1_3", "instance1_4"}
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

	cpuSysAll := window.New(W_SIZE, W_MULT)
	cpuSysAll.PushBack(float64(15000000))
	cpuTot1_1 := window.New(W_SIZE, W_MULT)
	cpuTot1_1.PushBack(float64(20000))
	cpuTot1_2 := window.New(W_SIZE, W_MULT)
	cpuTot1_2.PushBack(float64(60000))
	cpuTot1_3 := window.New(W_SIZE, W_MULT)
	cpuTot1_3.PushBack(float64(60000))
	cpuTot1_4 := window.New(W_SIZE, W_MULT)
	cpuTot1_4.PushBack(float64(70000))
	cpuTot2_1 := window.New(W_SIZE, W_MULT)
	cpuTot2_1.PushBack(float64(40000))

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
		cpuSysAll,
	}
	cpu1_4 := CpuStats{
		cpuTot1_4,
		cpuSysAll,
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

	system := SystemStats{15000000}

	mockStats := GruStats{
		Service:  services,
		Instance: instances,
		System:   system,
	}

	return mockStats
}
