package policy

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/data"
	"github.com/elleFlorio/gru/enum"
	res "github.com/elleFlorio/gru/resources"
)

const c_EPSILON = 0.09

func init() {
	cfg.SetServices(createServices())
	cfg.SetNode(createNode())
	cfg.SetPolicy(createMockPolicy())

	res.GetResources().CPU.Total = 4
	res.GetResources().Memory.Total = 4 * 1024 * 1024 * 1024
}

const c_SWAP_NAME = "swap"
const c_SCALEIN_NAME = "scalein"
const c_SCALEOUT_NAME = "scaleout"

func TestList(t *testing.T) {
	list := List()
	assert.Contains(t, list, c_SWAP_NAME)
	assert.Contains(t, list, c_SCALEIN_NAME)
	assert.Contains(t, list, c_SCALEOUT_NAME)
}

func TestListPolicyActions(t *testing.T) {
	actSwap := ListPolicyActions(c_SWAP_NAME)
	assert.Contains(t, actSwap, "stop")
	assert.Contains(t, actSwap, "start")

	actScalein := ListPolicyActions(c_SCALEIN_NAME)
	assert.Contains(t, actScalein, "stop")

	actScaleout := ListPolicyActions(c_SCALEOUT_NAME)
	assert.Contains(t, actScaleout, "start")
}

func TestScaleIn(t *testing.T) {
	shared := createSharedData()
	creator := &scaleinCreator{}

	w1 := creator.computeWeight("service1", shared)
	w2 := creator.computeWeight("service2", shared)
	w3 := creator.computeWeight("service3", shared)
	assert.InDelta(t, 0.0, w1, c_EPSILON)
	assert.InDelta(t, 0.16, w2, c_EPSILON)
	assert.InDelta(t, 0.0, w3, c_EPSILON)
}

func TestScaleOut(t *testing.T) {
	shared := createSharedData()
	creator := &scaleoutCreator{}

	w1 := creator.computeWeight("service1", shared)
	w2 := creator.computeWeight("service2", shared)

	res.GetResources().CPU.Used = 4
	w3 := creator.computeWeight("service3", shared)
	res.GetResources().CPU.Used = 0

	w4 := creator.computeWeight("service4", shared)
	assert.InDelta(t, 0.5, w1, c_EPSILON)
	assert.InDelta(t, 0.0, w2, c_EPSILON)
	assert.InDelta(t, 0.0, w3, c_EPSILON)
	assert.InDelta(t, 0.25, w4, c_EPSILON)
}

func TestSwap(t *testing.T) {
	shared := createSharedData()
	creator := &swapCreator{}

	srvList := []string{
		"service1",
		"service2",
		"service3",
		"service4",
		"service5",
	}

	pairs := creator.createSwapPairs(srvList)
	assert.Len(t, pairs, 2)

	pair1 := pairs["service1"]
	pair2 := pairs["service2"]
	assert.Len(t, pair1, 3)
	assert.Len(t, pair2, 3)
	assert.Contains(t, pair1, "service3")
	assert.Contains(t, pair1, "service5")
	assert.Contains(t, pair2, "service3")
	assert.Contains(t, pair2, "service5")

	res.GetResources().CPU.Used = 4
	w13 := creator.computeWeight("service1", "service3", shared)
	res.GetResources().CPU.Used = 0

	w14 := creator.computeWeight("service1", "service4", shared)
	w15 := creator.computeWeight("service1", "service5", shared)
	assert.InDelta(t, 0.16, w13, c_EPSILON)
	assert.Equal(t, 0.0, w14)
	assert.Equal(t, 0.0, w15)

	res.GetResources().CPU.Used = 4
	w23 := creator.computeWeight("service2", "service3", shared)
	res.GetResources().CPU.Used = 0

	w24 := creator.computeWeight("service2", "service4", shared)
	w25 := creator.computeWeight("service2", "service5", shared)
	assert.InDelta(t, 0.8, w23, c_EPSILON)
	assert.Equal(t, 0.0, w24)
	assert.Equal(t, 0.0, w25)
}

func TestCreatePolicy(t *testing.T) {
	shared := createSharedData()
	srvList := []string{
		"service1",
		"service2",
		"service3",
		"service4",
		"service5",
	}

	policies := CreatePolicies(srvList, shared)
	assert.Len(t, policies, 17)

}

// ######## MOCK ########

func createServices() []cfg.Service {
	srv1 := cfg.Service{}
	srv1.Name = "service1"
	srv1.Instances.Running = []string{"instance1_1"}
	srv1.Docker.CPUnumber = 2
	srv1.Analytics = []string{"LOAD"}

	srv2 := cfg.Service{}
	srv2.Name = "service2"
	srv2.Instances.Running = []string{"instance2_1"}
	srv2.Instances.Pending = []string{"instance2_2"}
	srv2.Docker.CPUnumber = 2
	srv2.Analytics = []string{"LOAD"}

	srv3 := cfg.Service{}
	srv3.Name = "service3"
	srv3.Docker.CPUnumber = 2
	srv3.Analytics = []string{"LOAD"}

	srv4 := cfg.Service{}
	srv4.Name = "service4"
	srv4.Docker.CPUnumber = 1
	srv4.Analytics = []string{"LOAD"}

	srv5 := cfg.Service{}
	srv5.Name = "service5"
	srv5.Docker.CPUnumber = 2
	srv5.Analytics = []string{"LOAD"}

	services := []cfg.Service{srv1, srv2, srv3, srv4, srv5}

	return services
}

func createSharedData() data.Shared {
	shared := data.Shared{}

	srvData1 := data.SharedData{
		BaseShared: map[string]float64{
			enum.METRIC_CPU_AVG.ToString(): 0.5,
		},
		UserShared: map[string]float64{
			"LOAD": 1.0,
		},
	}
	srvShared1 := data.ServiceShared{
		Data:   srvData1,
		Active: true,
	}

	srvData2 := data.SharedData{
		BaseShared: map[string]float64{
			enum.METRIC_CPU_AVG.ToString(): 0.2,
		},
		UserShared: map[string]float64{
			"LOAD": 0.5,
		},
	}
	srvShared2 := data.ServiceShared{
		Data:   srvData2,
		Active: true,
	}

	srvData3 := data.SharedData{
		BaseShared: map[string]float64{
			enum.METRIC_CPU_AVG.ToString(): 0.8,
		},
		UserShared: map[string]float64{
			"LOAD": 0.9,
		},
	}
	srvShared3 := data.ServiceShared{
		Data:   srvData3,
		Active: false,
	}

	srvData4 := data.SharedData{
		BaseShared: map[string]float64{
			enum.METRIC_CPU_AVG.ToString(): 0.8,
		},
		UserShared: map[string]float64{
			"LOAD": 0.9,
		},
	}
	srvShared4 := data.ServiceShared{
		Data:   srvData4,
		Active: false,
	}

	srvData5 := data.SharedData{
		BaseShared: map[string]float64{
			enum.METRIC_CPU_AVG.ToString(): 0.0,
		},
		UserShared: map[string]float64{
			"LOAD": 0.0,
		},
	}
	srvShared5 := data.ServiceShared{
		Data:   srvData5,
		Active: false,
	}

	shared.Service = map[string]data.ServiceShared{
		"service1": srvShared1,
		"service2": srvShared2,
		"service3": srvShared3,
		"service4": srvShared4,
		"service5": srvShared5,
	}

	return shared
}

func createNode() cfg.Node {
	n := cfg.Node{}
	n.Constraints.BaseServices = []string{"service3"}
	return n
}

func createMockPolicy() cfg.Policy {
	policy := cfg.Policy{
		Scalein: cfg.PolicyConfig{
			Enable:    true,
			Threshold: 0.3,
		},
		Scaleout: cfg.PolicyConfig{
			Enable:    true,
			Threshold: 0.8,
		},
		Swap: cfg.PolicyConfig{
			Enable:    true,
			Threshold: 0.6,
		},
	}

	return policy
}
