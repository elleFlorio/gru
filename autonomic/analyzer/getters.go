package analyzer

import (
	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/storage"
)

func getAnalytics() GruAnalytics {
	data, _ := storage.GetClusterData(enum.ANALYTICS)
	analytics, _ := convertDataToAnalytics(data)
	return analytics
}

func GetNodeAnalytics() GruAnalytics {
	return getAnalytics()
}

func GetServiceAnalytics(name string) ServiceAnalytics {
	return getAnalytics().Service[name]
}

func GetServicesAnalytics() map[string]ServiceAnalytics {
	return getAnalytics().Service
}

func GetSystemAnalytics() SystemAnalytics {
	return getAnalytics().System
}

func GetClusterAnalytics() ClusterAnalytics {
	return getAnalytics().Cluster
}
