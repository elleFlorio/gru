package analyzer

import (
	"github.com/elleFlorio/gru/autonomic/monitor"
)

func CreateMockAnalytics() GruAnalytics {
	mockAnalytics := GruAnalytics{
		Service:  make(map[string]ServiceAnalytics),
		Instance: make(map[string]InstanceAnalytics),
	}
	mockStats := monitor.CreateMockStats()

	for _, name := range monitor.ListMockServices() {
		updateInstances(name, &mockAnalytics, &mockStats, monitor.MaxNumberOfEntryInHistory())
		mockCpuTot, mockCpuAvg := computeServiceCpuPerc(name, &mockAnalytics, &mockStats)
		mockSrv := mockAnalytics.Service[name]
		mockSrv.CpuTot = mockCpuTot
		mockSrv.CpuAvg = mockCpuAvg
		mockAnalytics.Service[name] = mockSrv
		updateSystemInstances(&mockAnalytics)
	}

	return mockAnalytics
}
