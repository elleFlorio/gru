package monitor

func GetNodeStats() GruStats {
	services := GetServicesStats()
	instances := GetInstancesStats()
	system := GetSystemStats()

	return GruStats{
		services,
		instances,
		system,
	}
}

func GetServiceStats(name string) ServiceStats {
	return gruStats.Service[name]
}

func GetServicesStats() map[string]ServiceStats {
	return gruStats.Service
}

func GetInstanceStats(id string) InstanceStats {
	return gruStats.Instance[id]
}

func GetInstancesStats() map[string]InstanceStats {
	return gruStats.Instance
}

func GetSystemStats() SystemStats {
	return gruStats.System
}
