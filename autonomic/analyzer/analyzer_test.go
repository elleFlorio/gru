package analyzer

import (
	"testing"

	"github.com/elleFlorio/gru/autonomic/monitor"
	"github.com/stretchr/testify/assert"
)

func TestGetActiveInstances(t *testing.T) {
	mockStats := monitor.CreateMockStats()
	maxHist := monitor.MaxNumberOfEntryInHistory()
	running := mockStats.Service["service1"].Instances.Running

	active, pending := getActiveInstances(running, &mockStats, maxHist)
	assert.Equal(t, 2, len(active), "Number of active instances should be 2")
	assert.Equal(t, 1, len(pending), "Number of pending instances should be 1")

}

func TestUpdateInstances(t *testing.T) {
	mockAnalytics := &GruAnalytics{
		Service:  make(map[string]ServiceAnalytics),
		Instance: make(map[string]InstanceAnalytics),
	}
	mockStats := monitor.CreateMockStats()
	maxHist := monitor.MaxNumberOfEntryInHistory()

	for _, name := range monitor.ListMockServices() {
		updateInstances(name, mockAnalytics, &mockStats, maxHist)
	}

	assert.Contains(t, mockAnalytics.Service["service1"].Instances.Active, "instance1_1", "Active instances should contain instance1_1")
	assert.Contains(t, mockAnalytics.Service["service1"].Instances.Pending, "instance1_4", "Pending instances should contain instance1_4")
	assert.Contains(t, mockAnalytics.Service["service1"].Instances.Paused, "instance1_3", "Paused instances should contain instance1_3")
	assert.Contains(t, mockAnalytics.Service["service1"].Instances.Stopped, "instance1_0", "Paused instances should contain instance1_0")

}

// computeInstanceCpuPerc(instCpus []float64, sysCpus []float64) float64
func TestComputeInstanceCpuPerc(t *testing.T) {
	mockInstCpus := []float64{10000, 20000, 30000, 40000, 50000, 60000}
	mockSysCpus := []float64{1000000, 1100000, 1200000, 1300000, 1400000, 1500000}

	mockPerc := computeInstanceCpuPerc(mockInstCpus, mockSysCpus)

	assert.Equal(t, 0.1, mockPerc, "Computed percentage should be 1%")
}

// computeServiceCpuPerc(name string, analytics *GruAnalytics, stats *monitor.GruStats) float64
func TestComputeServiceCpuPerc(t *testing.T) {
	mockAnalytics := &GruAnalytics{
		Service:  make(map[string]ServiceAnalytics),
		Instance: make(map[string]InstanceAnalytics),
	}
	mockStats := monitor.CreateMockStats()

	for _, name := range monitor.ListMockServices() {
		updateInstances(name, mockAnalytics, &mockStats, monitor.MaxNumberOfEntryInHistory())
	}

	mockCpuTotSrv1, mockCpuAvgSrv1 := computeServiceCpuPerc("service1", mockAnalytics, &mockStats)
	mockCpuTotSrv2, mockCpuAvgSrv2 := computeServiceCpuPerc("service2", mockAnalytics, &mockStats)

	assert.Equal(t, 0.7, mockCpuTotSrv1, "Cpu total of service1 should be 0.7")
	assert.Equal(t, 0.4, mockCpuTotSrv2, "Cpu total of service2 should be 0.4")

	assert.Equal(t, 0.35, mockCpuAvgSrv1, "Cpu average of service1 should be 0.35")
	assert.Equal(t, 0.4, mockCpuAvgSrv2, "Cpu average of service2 should be 0.4")

}

func TestUpdateSysInstances(t *testing.T) {
	mockAnalytics := &GruAnalytics{
		Service:  make(map[string]ServiceAnalytics),
		Instance: make(map[string]InstanceAnalytics),
	}
	mockStats := monitor.CreateMockStats()
	mockAll := 0
	mockActive := 0
	mockPending := 0
	mockPaused := 0
	mockStopped := 0

	for _, name := range monitor.ListMockServices() {
		updateInstances(name, mockAnalytics, &mockStats, monitor.MaxNumberOfEntryInHistory())
		mockAll += len(mockAnalytics.Service[name].Instances.All)
		mockActive += len(mockAnalytics.Service[name].Instances.Active)
		mockPending += len(mockAnalytics.Service[name].Instances.Pending)
		mockPaused += len(mockAnalytics.Service[name].Instances.Paused)
		mockStopped += len(mockAnalytics.Service[name].Instances.Stopped)
	}

	updateSystemInstances(mockAnalytics)
	assert.Equal(t, mockAll, len(mockAnalytics.System.Instances.All),
		"Number of all instances in system should be", mockAll)
	assert.Equal(t, mockActive, len(mockAnalytics.System.Instances.Active),
		"Number of active instances in system should be", mockActive)
	assert.Equal(t, mockPending, len(mockAnalytics.System.Instances.Pending),
		"Number of pending instances in system should be", mockPending)
	assert.Equal(t, mockPaused, len(mockAnalytics.System.Instances.Paused),
		"Number of paused instances in system should be", mockPaused)
	assert.Equal(t, mockStopped, len(mockAnalytics.System.Instances.Stopped),
		"Number of stopped instances in system should be", mockStopped)
}
