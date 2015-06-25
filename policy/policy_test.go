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
		"scalein":  &ScaleIn{},
		"scaleout": &ScaleOut{},
	}
}

func TestWeight(t *testing.T) {
	mockServices := service.CreateMockServices()
	mockAnalytics := createMockAnalytics()
	delta := 0.0000001

	// ScaleDown
	weight0 := plc["scalein"].Weight(&mockServices[0], mockAnalytics)
	assert.Equal(t, 1.0, weight0, "scalein: weight0 should be 1")
	weight1 := plc["scalein"].Weight(&mockServices[1], mockAnalytics)
	assert.Equal(t, 0.5, weight1, "scalein: weight1 should be 0.5")
	weight2 := plc["scalein"].Weight(&mockServices[2], mockAnalytics)
	assert.Equal(t, 0.0, weight2, "scalein: weight2 should be 0")

	//ScaleUp
	weight0 = plc["scaleout"].Weight(&mockServices[0], mockAnalytics)
	assert.Equal(t, 0.0, weight0, "scaleout: weight0 should be 0")
	weight1 = plc["scaleout"].Weight(&mockServices[1], mockAnalytics)
	assert.Equal(t, 0.0, weight1, "scaleout: weight1 should be 0")
	weight2 = plc["scaleout"].Weight(&mockServices[2], mockAnalytics)
	assert.InDelta(t, 0.5, weight2, delta, "scaleout: weight2 should be 0")
}

// Too specific to use the one provided by the analyzer
func createMockAnalytics() *analyzer.GruAnalytics {
	active1 := []string{"instance1_1", "instance1_2", "instance1_3"}
	instances1 := analyzer.InstanceStatus{
		Active: active1,
	}
	service1 := analyzer.ServiceAnalytics{
		CpuAvg:    0,
		Instances: instances1}

	active2 := []string{"instance2_1", "instance_2_2"}
	instances2 := analyzer.InstanceStatus{
		Active: active2,
	}
	service2 := analyzer.ServiceAnalytics{
		CpuAvg:    0.2,
		Instances: instances2,
	}

	active3 := []string{"instance3_1", "instance3_2", "instance3_3", "instance3_4"}
	instances3 := analyzer.InstanceStatus{
		Active: active3,
	}
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

	return &mockAnalytics
}
