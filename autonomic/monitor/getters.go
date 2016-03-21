package monitor

import "github.com/elleFlorio/gru/data"

func GetNodeStats() data.GruStats {
	services := GetServicesStats()
	instances := GetInstancesStats()
	system := GetSystemStats()

	return data.GruStats{
		services,
		instances,
		system,
	}
}

func GetServiceStats(name string) data.ServiceStats {
	return gruStats.Service[name]
}

func GetServicesStats() map[string]data.ServiceStats {
	return gruStats.Service
}

func GetInstanceStats(id string) data.InstanceStats {
	return gruStats.Instance[id]
}

func GetInstancesStats() map[string]data.InstanceStats {
	return gruStats.Instance
}

func GetSystemStats() data.SystemStats {
	return gruStats.System
}
