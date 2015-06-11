package policy

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/autonomic/analyzer"
	"github.com/elleFlorio/gru/service"
)

var plc map[string]GruPolicy

func init() {
	plc = map[string]GruPolicy{
		"scaleDown": &ScaleDown{},
		"scaleUp":   &ScaleUp{},
	}
}

func TestWeight(t *testing.T) {
	mockServices := createMockServices()
	mockAnalytics := createMockAnalytics()
	delta := 0.0000001

	// ScaleDown
	weight0 := plc["scaleDown"].Weight(&mockServices[0], &mockAnalytics)
	assert.Equal(t, 1.0, weight0, "scaleDown: weight0 should be 1")
	weight1 := plc["scaleDown"].Weight(&mockServices[1], &mockAnalytics)
	assert.Equal(t, 0.5, weight1, "scaleDown: weight1 should be 0.5")
	weight2 := plc["scaleDown"].Weight(&mockServices[2], &mockAnalytics)
	assert.Equal(t, 0.0, weight2, "scaleDown: weight2 should be 0")

	//ScaleUp
	weight0 = plc["scaleUp"].Weight(&mockServices[0], &mockAnalytics)
	assert.Equal(t, 0.0, weight0, "scaleUp: weight0 should be 0")
	weight1 = plc["scaleUp"].Weight(&mockServices[1], &mockAnalytics)
	assert.Equal(t, 0.0, weight1, "scaleUp: weight1 should be 0")
	weight2 = plc["scaleUp"].Weight(&mockServices[2], &mockAnalytics)
	// FIXME is there a way to avoid this step???
	if weight2 < 0.5+delta && weight2 > 0.5-delta {
		weight2 = 0.5
	}
	assert.Equal(t, 0.5, weight2, "scaleUp: weight2 should be 0")
}

func createMockServices() []service.Service {
	mockService1 := service.Service{
		Name: "service1",
		Constraints: service.Constraints{
			CpuMin:    0.4,
			CpuMax:    0.6,
			MinActive: 1,
			MaxActive: 3,
		},
	}

	mockService2 := service.Service{
		Name: "service2",
		Constraints: service.Constraints{
			CpuMin:    0.4,
			CpuMax:    0.7,
			MinActive: 1,
			MaxActive: 2,
		},
	}

	mockService3 := service.Service{
		Name: "service3",
		Constraints: service.Constraints{
			CpuMin:    0.4,
			CpuMax:    0.6,
			MinActive: 2,
			MaxActive: 5,
		},
	}

	return []service.Service{
		mockService1,
		mockService2,
		mockService3,
	}
}

func createMockAnalytics() analyzer.GruAnalytics {
	instances1 := []string{"instance1_1", "instance1_2", "instance1_3"}
	service1 := analyzer.ServiceAnalytics{
		CpuAvg:    0,
		Instances: instances1}

	instances2 := []string{"instance2_1", "instance_2_2"}
	service2 := analyzer.ServiceAnalytics{
		CpuAvg:    0.2,
		Instances: instances2,
	}

	instances3 := []string{"instance3_1", "instance3_2", "instance3_3", "instance3_4"}
	service3 := analyzer.ServiceAnalytics{
		CpuAvg:    0.8,
		Instances: instances3,
	}

	services := map[string]analyzer.ServiceAnalytics{
		"service1": service1,
		"service2": service2,
		"service3": service3,
	}

	mockAnalytics := analyzer.GruAnalytics{
		Service: services,
	}

	return mockAnalytics
}
