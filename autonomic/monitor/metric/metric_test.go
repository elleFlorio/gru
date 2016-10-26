package metric

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/data"
	"github.com/elleFlorio/gru/enum"
	res "github.com/elleFlorio/gru/resources"
	srv "github.com/elleFlorio/gru/service"
	"github.com/elleFlorio/gru/utils"
)

func init() {
	res.CreateMockResources(1, "4GB", 0, "0g")
	resetMockServices()
	Initialize(srv.List())
}

func TestAddInstance(t *testing.T) {
	id := "id1"
	defer delete(instancesMetrics, id)

	AddInstance(id)
	_, ok := instancesMetrics[id]
	assert.True(t, ok)
}

func TestRemoveInstance(t *testing.T) {
	id := "id1"
	instancesMetrics[id] = Metric{
		BaseMetrics: make(map[string][]float64),
	}

	RemoveInstance(id)
	assert.Empty(t, instancesMetrics)
}

func TestUpdateCpuMetric(t *testing.T) {
	id := "id1"
	defer delete(instancesMetrics, id)

	instancesMetrics[id] = Metric{
		BaseMetrics: make(map[string][]float64),
	}
	toAddInst := []float64{1, 2, 3, 4, 5}
	toAddSys := []float64{1, 2, 3, 4, 5}

	UpdateCpuMetric(id, toAddInst, toAddSys)
	cpuInst := instancesMetrics[id].BaseMetrics[enum.METRIC_CPU_INST.ToString()]
	cpuSys := instancesMetrics[id].BaseMetrics[enum.METRIC_CPU_SYS.ToString()]
	assert.Equal(t, toAddInst, cpuInst)
	assert.Equal(t, toAddSys, cpuSys)

	// check logs for error
	UpdateCpuMetric("pippo", toAddInst, toAddSys)
}

func TestUpdateMemMetric(t *testing.T) {
	id := "id1"
	defer delete(instancesMetrics, id)

	instancesMetrics[id] = Metric{
		BaseMetrics: make(map[string][]float64),
	}
	toAdd := []float64{1, 2, 3, 4, 5}

	UpdateMemMetric(id, toAdd)
	mem := instancesMetrics[id].BaseMetrics[enum.METRIC_MEM_INST.ToString()]
	assert.Equal(t, toAdd, mem)

	// check logs for errors
	UpdateMemMetric("pippo", toAdd)
}

func TestUpdateUserMetric(t *testing.T) {
	defer resetMetrics()

	service := "service1"
	metric := "response_time"
	toAdd := []float64{1, 2, 3, 4, 5}

	UpdateUserMetric(service, metric, toAdd)
	assert.Equal(t, toAdd, servicesMetrics[service].UserMetrics[metric])

	// check logs for errors
	UpdateUserMetric("pippo", metric, toAdd)

}

func TestIsReadyForRunning(t *testing.T) {
	id := "id1"
	defer delete(instancesMetrics, id)

	instancesMetrics[id] = Metric{
		BaseMetrics: make(map[string][]float64),
	}

	values1 := []float64{1, 2}
	values2 := []float64{1, 2, 3, 4}
	values3 := []float64{1, 2, 3, 4, 5}
	instancesMetrics[id].BaseMetrics[enum.METRIC_CPU_INST.ToString()] = values1
	instancesMetrics[id].BaseMetrics[enum.METRIC_CPU_SYS.ToString()] = values2
	instancesMetrics[id].BaseMetrics[enum.METRIC_MEM_INST.ToString()] = values3
	thr1 := len(values1)
	thr2 := len(values2)
	thr3 := len(values3)

	assert.True(t, IsReadyForRunning(id, thr1))
	assert.False(t, IsReadyForRunning(id, thr2))
	assert.False(t, IsReadyForRunning(id, thr3))
}

func TestComputeInstanceCpuPerc(t *testing.T) {
	mockInstCpus := []float64{10000, 20000, 30000, 40000, 50000, 60000}
	mockSysCpus := []float64{1000000, 1100000, 1200000, 1300000, 1400000, 1500000}

	mockPerc := computeInstanceCpuPerc(mockInstCpus, mockSysCpus)
	assert.Equal(t, 0.1, mockPerc)

	mockInstCpus = []float64{10000, 10000, 10000, 10000, 10000, 10000}

	mockPerc = computeInstanceCpuPerc(mockInstCpus, mockSysCpus)
	assert.Equal(t, 0.0, mockPerc)
}

func TestComputeInstancesMetrics(t *testing.T) {
	id1 := "instance1_1"
	id2 := "instance2_1"
	instMemBytes, _ := utils.RAMInBytes("512m")
	instMemValue := float64(instMemBytes)
	defer delete(instancesMetrics, id1)
	defer delete(instancesMetrics, id2)

	instancesMetrics[id1] = Metric{
		BaseMetrics: make(map[string][]float64),
	}
	instancesMetrics[id1].BaseMetrics[enum.METRIC_CPU_INST.ToString()] = []float64{10000, 20000, 30000, 40000, 50000, 60000}
	instancesMetrics[id1].BaseMetrics[enum.METRIC_CPU_SYS.ToString()] = []float64{1000000, 1100000, 1200000, 1300000, 1400000, 1500000}
	instancesMetrics[id1].BaseMetrics[enum.METRIC_MEM_INST.ToString()] = []float64{instMemValue, instMemValue}

	instancesMetrics[id2] = Metric{
		BaseMetrics: make(map[string][]float64),
	}
	instancesMetrics[id2].BaseMetrics[enum.METRIC_CPU_INST.ToString()] = []float64{10000, 30000, 50000, 70000, 110000, 130000}
	instancesMetrics[id2].BaseMetrics[enum.METRIC_CPU_SYS.ToString()] = []float64{1000000, 1100000, 1200000, 1300000, 1400000, 1500000}
	instancesMetrics[id2].BaseMetrics[enum.METRIC_MEM_INST.ToString()] = []float64{instMemValue, instMemValue}

	instMet := computeInstancesMetrics()
	assert.Equal(t, 0.1, instMet[id1].BaseMetrics[enum.METRIC_CPU_AVG.ToString()])
	assert.Equal(t, 0.24, instMet[id2].BaseMetrics[enum.METRIC_CPU_AVG.ToString()])
	assert.Equal(t, 0.5, instMet[id1].BaseMetrics[enum.METRIC_MEM_AVG.ToString()])
}

func TestComputeServicesMetrics(t *testing.T) {
	defer resetMetrics()

	instMetrics := make(map[string]data.MetricData)
	instMetrics["instance1_1"] = data.MetricData{
		BaseMetrics: map[string]float64{
			enum.METRIC_CPU_AVG.ToString(): 0.4,
			enum.METRIC_MEM_AVG.ToString(): 0.4,
		},
	}
	instMetrics["instance1_2"] = data.MetricData{
		BaseMetrics: map[string]float64{
			enum.METRIC_CPU_AVG.ToString(): 0.6,
			enum.METRIC_MEM_AVG.ToString(): 0.6,
		},
	}
	instMetrics["instance2_1"] = data.MetricData{
		BaseMetrics: map[string]float64{
			enum.METRIC_CPU_AVG.ToString(): 0.5,
			enum.METRIC_MEM_AVG.ToString(): 0.5,
		},
	}

	servicesMetrics["service1"].UserMetrics["response_time"] = []float64{1000, 2000, 3000, 4000, 5000}
	servicesMetrics["service2"].UserMetrics["response_time"] = []float64{1000, 2000, 3000, 4000, 5000}

	serviceMet := computeServicesMetrics(instMetrics)
	assert.NotEmpty(t, serviceMet)
	assert.InEpsilon(t, 0.5, serviceMet["service1"].BaseMetrics[enum.METRIC_CPU_AVG.ToString()], 0.0001)
	assert.InEpsilon(t, 0.5, serviceMet["service2"].BaseMetrics[enum.METRIC_CPU_AVG.ToString()], 0.0001)
	assert.Equal(t, 0.0, serviceMet["service3"].BaseMetrics[enum.METRIC_CPU_AVG.ToString()])
	assert.Equal(t, 0.0, serviceMet["service3"].BaseMetrics[enum.METRIC_MEM_AVG.ToString()])
	assert.Equal(t, 3000.0, serviceMet["service1"].UserMetrics["response_time"])
	assert.Equal(t, 3000.0, serviceMet["service2"].UserMetrics["response_time"])
}

func TestComputeSysMetrics(t *testing.T) {
	defer resetMetrics()

	instMetrics := make(map[string]data.MetricData)
	instMetrics["instance1_1"] = data.MetricData{
		BaseMetrics: map[string]float64{
			enum.METRIC_CPU_AVG.ToString(): 0.4,
		},
	}
	instMetrics["instance1_2"] = data.MetricData{
		BaseMetrics: map[string]float64{
			enum.METRIC_CPU_AVG.ToString(): 0.6,
		},
	}
	instMetrics["instance2_1"] = data.MetricData{
		BaseMetrics: map[string]float64{
			enum.METRIC_CPU_AVG.ToString(): 0.8,
		},
	}

	res.CreateMockResources(4, "1G", 0, "512M")
	sysMet := computeSysMetrics(instMetrics)
	assert.NotEmpty(t, sysMet)
	assert.InEpsilon(t, 0.65, sysMet.BaseMetrics[enum.METRIC_CPU_AVG.ToString()], 0.0001)
	assert.Equal(t, 0.5, sysMet.BaseMetrics[enum.METRIC_MEM_AVG.ToString()])
}

func TestResetMetrics(t *testing.T) {
	id := "id1"
	defer delete(instancesMetrics, id)

	instancesMetrics[id] = Metric{
		BaseMetrics: make(map[string][]float64),
	}
	instancesMetrics[id].BaseMetrics[enum.METRIC_CPU_INST.ToString()] = []float64{10000, 20000, 30000, 40000, 50000, 60000}
	instancesMetrics[id].BaseMetrics[enum.METRIC_CPU_SYS.ToString()] = []float64{1000000, 1100000, 1200000, 1300000, 1400000, 1500000}
	servicesMetrics["service1"].UserMetrics["response_time"] = []float64{1000, 2000, 3000, 4000, 5000}
	servicesMetrics["service2"].UserMetrics["response_time"] = []float64{1000, 2000, 3000, 4000, 5000}

	resetMetrics()
	assert.Len(t, servicesMetrics, len(srv.List()))
	assert.Empty(t, servicesMetrics["service1"].UserMetrics["response_time"])
	assert.Empty(t, servicesMetrics["service2"].UserMetrics["response_time"])
	assert.Empty(t, servicesMetrics["service3"].UserMetrics["response_time"])
	assert.Len(t, instancesMetrics, 1)
	assert.Empty(t, instancesMetrics[id].BaseMetrics[enum.METRIC_CPU_INST.ToString()])
	assert.Empty(t, instancesMetrics[id].BaseMetrics[enum.METRIC_CPU_SYS.ToString()])

}

func resetMockServices() {
	mockServices := srv.CreateMockServices()
	cfg.SetServices(mockServices)
}
