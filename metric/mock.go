package metric

import (
	"github.com/elleFlorio/gru/service"
)

func CreateMetricsMockConfig(serviceName string) map[string]interface{} {
	config := make(map[string]interface{})
	switch serviceName {
	case "influxdb":
		config = createInfluxMockConfig()
	}

	return config

}

func createInfluxMockConfig() map[string]interface{} {
	return map[string]interface{}{
		"Url":      "http://mockip:mockport",
		"DbName":   "mockDB",
		"Username": "mockUser",
		"Password": "mockPwd",
	}
}

func CreateMockMetrics() GruMetric {
	service.UpdateServices(service.CreateMockServices())
	mockMetrics := newMetrics()
	return mockMetrics
}
