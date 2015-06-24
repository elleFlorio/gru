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

	active, pending := getActiveInstances(running, mockStats, maxHist)
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
		updateInstances(name, mockAnalytics, mockStats, maxHist)
	}

	assert.Contains(t, mockAnalytics.Service["service1"].Instances.Active, "instance1_1", "Active instances should contain instance1_1")
	assert.Contains(t, mockAnalytics.Service["service1"].Instances.Pending, "instance1_4", "Pending instances should contain instance1_4")
	assert.Contains(t, mockAnalytics.Service["service1"].Instances.Paused, "instance1_3", "Paused instances should contain instance1_3")
	assert.Contains(t, mockAnalytics.Service["service1"].Instances.Stopped, "instance1_0", "Paused instances should contain instance1_0")

}

// computeInstanceCpuPerc(instCpus []float64, sysCpus []float64) float64
func TestComputeInstanceCpuPerc(t *testing.T) {
	mockInstCpus := []float64{10000, 20000, 30000, 40000, 50000, 60000}
	mockSysCpus := []float64{100000000, 110000000, 120000000, 130000000, 140000000, 150000000}

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
		updateInstances(name, mockAnalytics, mockStats, monitor.MaxNumberOfEntryInHistory())
	}

	mockCpuSrv1 := computeServiceCpuPerc("service1", mockAnalytics, mockStats)
	mockCpuSrv2 := computeServiceCpuPerc("service2", mockAnalytics, mockStats)

	assert.Equal(t, 0.35, mockCpuSrv1, "Cpu average of service1 should be 0.35")
	assert.Equal(t, 0.4, mockCpuSrv2, "Cpu average of service1 should be 0.4")

}
