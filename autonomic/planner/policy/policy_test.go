package policy

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/autonomic/analyzer"
	cfg "github.com/elleFlorio/gru/configuration"
	res "github.com/elleFlorio/gru/resources"
)

const c_EPSILON = 0.09

func init() {
	cfg.SetServices(createServices())
	cfg.SetNode(createNode())

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
	analytics := createAnalytics()
	creator := &scaleinCreator{}

	w1 := creator.computeWeight("service1", analytics)
	w2 := creator.computeWeight("service2", analytics)
	w3 := creator.computeWeight("service3", analytics)
	assert.InDelta(t, 0.0, w1, c_EPSILON)
	assert.InDelta(t, 0.16, w2, c_EPSILON)
	assert.InDelta(t, 0.0, w3, c_EPSILON)
}

func TestScaleOut(t *testing.T) {
	analytics := createAnalytics()
	creator := &scaleoutCreator{}

	w1 := creator.computeWeight("service1", analytics)
	w2 := creator.computeWeight("service2", analytics)
	w3 := creator.computeWeight("service3", analytics)
	w4 := creator.computeWeight("service4", analytics)
	assert.InDelta(t, 0.5, w1, c_EPSILON)
	assert.InDelta(t, 0.0, w2, c_EPSILON)
	assert.InDelta(t, 0.0, w3, c_EPSILON)
	assert.InDelta(t, 0.25, w4, c_EPSILON)
}

func TestSwap(t *testing.T) {
	analytics := createAnalytics()
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

	w13 := creator.computeWeight("service1", "service3", analytics)
	w14 := creator.computeWeight("service1", "service4", analytics)
	w15 := creator.computeWeight("service1", "service5", analytics)
	assert.InDelta(t, 0.16, w13, c_EPSILON)
	assert.Equal(t, 0.0, w14)
	assert.Equal(t, 0.0, w15)

	w23 := creator.computeWeight("service2", "service3", analytics)
	w24 := creator.computeWeight("service2", "service4", analytics)
	w25 := creator.computeWeight("service2", "service5", analytics)
	assert.InDelta(t, 0.8, w23, c_EPSILON)
	assert.Equal(t, 0.0, w24)
	assert.Equal(t, 0.0, w25)
}

func TestCreatePolicy(t *testing.T) {
	analytics := createAnalytics()

	srvList := []string{
		"service1",
		"service2",
		"service3",
		"service4",
		"service5",
	}

	policies := CreatePolicies(srvList, analytics)
	assert.Len(t, policies, 17)

}

// ######## MOCK ########

func createServices() []cfg.Service {
	srv1 := cfg.Service{}
	srv1.Name = "service1"
	srv1.Instances.Running = []string{"instance1_1"}

	srv2 := cfg.Service{}
	srv2.Name = "service2"
	srv2.Instances.Running = []string{"instance2_1"}
	srv2.Instances.Pending = []string{"instance2_2"}

	srv3 := cfg.Service{}
	srv3.Name = "service3"

	srv4 := cfg.Service{}
	srv4.Name = "service4"

	srv5 := cfg.Service{}
	srv5.Name = "service5"

	services := []cfg.Service{srv1, srv2, srv3, srv4, srv5}

	return services
}

func createAnalytics() analyzer.GruAnalytics {
	analytics := analyzer.GruAnalytics{}

	srv1A := analyzer.ServiceAnalytics{}
	srv1A.Load = 1.0
	srv1A.Resources.Cpu = 0.5
	srv1A.Resources.Available = 1.0

	srv2A := analyzer.ServiceAnalytics{}
	srv2A.Load = 0.5
	srv2A.Resources.Cpu = 0.2
	srv2A.Resources.Available = 1.0

	srv3A := analyzer.ServiceAnalytics{}
	srv3A.Load = 0.9
	srv3A.Resources.Cpu = 0.8
	srv3A.Resources.Available = 0.0

	srv4A := analyzer.ServiceAnalytics{}
	srv4A.Load = 0.9
	srv4A.Resources.Cpu = 0.8
	srv4A.Resources.Available = 1.0

	srv5A := analyzer.ServiceAnalytics{}
	srv5A.Resources.Available = 0.0

	analytics.Service = map[string]analyzer.ServiceAnalytics{
		"service1": srv1A,
		"service2": srv2A,
		"service3": srv3A,
		"service4": srv4A,
		"service5": srv5A,
	}

	return analytics
}

func createNode() cfg.Node {
	n := cfg.Node{}
	n.Constraints.BaseServices = []string{"service3"}
	return n
}
