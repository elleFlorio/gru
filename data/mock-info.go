package data

func CreateMockInfo() GruInfo {
	info := GruInfo{
		Service: make(map[string]ServiceInfo),
	}

	service1 := "service1"
	service2 := "service2"

	load1 := 0.7
	cpu1 := 0.9
	mem1 := 0.6
	res1 := 0.0
	act1 := true
	srvInfo1 := ServiceInfo{
		Load:      load1,
		Cpu:       cpu1,
		Memory:    mem1,
		Resources: res1,
		Active:    act1,
	}

	load2 := 0.4
	cpu2 := 0.5
	mem2 := 0.1
	res2 := 1.0
	act2 := true
	srvInfo2 := ServiceInfo{
		Load:      load2,
		Cpu:       cpu2,
		Memory:    mem2,
		Resources: res2,
		Active:    act2,
	}

	info.Service[service1] = srvInfo1
	info.Service[service2] = srvInfo2

	cpuSys := 0.8
	memSys := 0.7
	healthSys := 0.7
	activeSys := []string{service1}
	sys := SystemInfo{
		Cpu:            cpuSys,
		Memory:         memSys,
		Health:         healthSys,
		ActiveServices: activeSys,
	}

	info.System = sys

	return info

}
