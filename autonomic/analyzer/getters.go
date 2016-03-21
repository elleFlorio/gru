package analyzer

import (
	"github.com/elleFlorio/gru/data"
)

func getAnalytics() data.GruAnalytics {
	analytics, _ := data.GetAnalytics()

	return analytics
}

func GetNodeAnalytics() data.GruAnalytics {
	return getAnalytics()
}

func GetServiceAnalytics(name string) data.ServiceAnalytics {
	return getAnalytics().Service[name]
}

func GetServicesAnalytics() map[string]data.ServiceAnalytics {
	return getAnalytics().Service
}

func GetSystemAnalytics() data.SystemAnalytics {
	return getAnalytics().System
}

func GetClusterAnalytics() data.ClusterAnalytics {
	return getAnalytics().Cluster
}
