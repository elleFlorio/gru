package analyzer

import (
	"github.com/elleFlorio/gru/data"
)

func getAnalytics() data.GruAnalytics {
	analytics, _ := data.GetAnalytics()

	return analytics
}

func GetAnalytics() data.GruAnalytics {
	return getAnalytics()
}

func GetServiceAnalytics(name string) data.AnalyticData {
	return getAnalytics().Service[name]
}

func GetServicesAnalytics() map[string]data.AnalyticData {
	return getAnalytics().Service
}

func GetSystemAnalytics() data.AnalyticData {
	return getAnalytics().System
}
