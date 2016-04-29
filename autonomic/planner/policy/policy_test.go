package policy

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/data"
	res "github.com/elleFlorio/gru/resources"
)

const c_EPSILON = 0.09

func init() {
	cfg.SetServices(createServices())
	cfg.SetNode(createNode())
	cfg.SetTuning(createTuning())

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

	srv2 := cfg.Service{}
	srv2.Name = "service2"
	srv2.Instances.Running = []string{"instance2_1"}
	srv2.Instances.Pending = []string{"instance2_2"}
	srv2.Docker.CPUnumber = 2

	srv3 := cfg.Service{}
	srv3.Name = "service3"
	srv3.Docker.CPUnumber = 2

	srv4 := cfg.Service{}
	srv4.Name = "service4"
	srv4.Docker.CPUnumber = 1

	srv5 := cfg.Service{}
	srv5.Name = "service5"
	srv5.Docker.CPUnumber = 2

	services := []cfg.Service{srv1, srv2, srv3, srv4, srv5}

	return services
}

func createAnalytics() data.GruAnalytics {
	analytics := data.GruAnalytics{}

	srv1A := data.ServiceAnalytics{}
	srv1A.Load = 1.0
	srv1A.Resources.Cpu = 0.5
	srv1A.Resources.Available = 1.0

	srv2A := data.ServiceAnalytics{}
	srv2A.Load = 0.5
	srv2A.Resources.Cpu = 0.2
	srv2A.Resources.Available = 1.0

	srv3A := data.ServiceAnalytics{}
	srv3A.Load = 0.9
	srv3A.Resources.Cpu = 0.8
	srv3A.Resources.Available = 0.0

	srv4A := data.ServiceAnalytics{}
	srv4A.Load = 0.9
	srv4A.Resources.Cpu = 0.8
	srv4A.Resources.Available = 1.0

	srv5A := data.ServiceAnalytics{}
	srv5A.Resources.Available = 0.0

	analytics.Service = map[string]data.ServiceAnalytics{
		"service1": srv1A,
		"service2": srv2A,
		"service3": srv3A,
		"service4": srv4A,
		"service5": srv5A,
	}

	return analytics
}

func createSharedData() data.Shared {
	shared := data.Shared{}

	srv1A := data.ServiceShared{}
	srv1A.Load = 1.0
	srv1A.Cpu = 0.5
	srv1A.Resources = 1.0

	srv2A := data.ServiceShared{}
	srv2A.Load = 0.5
	srv2A.Cpu = 0.2
	srv2A.Resources = 1.0

	srv3A := data.ServiceShared{}
	srv3A.Load = 0.9
	srv3A.Cpu = 0.8
	srv3A.Resources = 0.0

	srv4A := data.ServiceShared{}
	srv4A.Load = 0.9
	srv4A.Cpu = 0.8
	srv4A.Resources = 1.0

	srv5A := data.ServiceShared{}
	srv5A.Resources = 0.0

	shared.Service = map[string]data.ServiceShared{
		"service1": srv1A,
		"service2": srv2A,
		"service3": srv3A,
		"service4": srv4A,
		"service5": srv5A,
	}

	return shared
}

func createNode() cfg.Node {
	n := cfg.Node{}
	n.Constraints.BaseServices = []string{"service3"}
	return n
}

func createTuning() cfg.Tuning {
	t := cfg.Tuning{}
	t.Policy.Scalein.Cpu = 0.3
	t.Policy.Scalein.Load = 0.3
	t.Policy.Scaleout.Cpu = 0.8
	t.Policy.Scaleout.Load = 0.8
	t.Policy.Swap.Cpu = 0.6
	t.Policy.Swap.Load = 0.6

	return t
}
