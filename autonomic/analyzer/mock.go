package analyzer

import (
	"github.com/elleFlorio/gru/autonomic/monitor"
	"github.com/elleFlorio/gru/service"
)

//FIXME when Analyzer is ready update mock generation
func CreateMockAnalytics() GruAnalytics {
	mockAnalytics := GruAnalytics{
		Service: make(map[string]ServiceAnalytics),
	}
	mockServices := service.CreateMockServices()
	service.UpdateServices(mockServices)
	mockStats := monitor.CreateMockStats()

	analyzeServices(&mockAnalytics, mockStats)
	analyzeSystem(&mockAnalytics, mockStats)
	computeNodeHealth(&mockAnalytics)

	return mockAnalytics
}
