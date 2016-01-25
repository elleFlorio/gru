package policy

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/autonomic/analyzer"
	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/enum"
	res "github.com/elleFlorio/gru/resources"
)

var plc map[string]GruPolicy

const c_EPSILON = 0.09

func init() {
	plc = map[string]GruPolicy{
		"scalein":  &ScaleIn{},
		"scaleout": &ScaleOut{},
	}

	res.GetResources().CPU.Total = 4
	res.GetResources().Memory.Total = 4 * 1024 * 1024 * 1024
}

func TestGetPolicies(t *testing.T) {
	pls := GetPolicies()
	names := make([]string, 0)
	actions := make([]enum.Actions, 0)
	for _, item := range pls {
		names = append(names, item.Name())
		actions = append(actions, item.Actions())
	}

	assert.Equal(t, len(plc), len(pls))
	for _, item := range plc {
		assert.Contains(t, names, item.Name())
		assert.Contains(t, actions, item.Actions())
	}
}

func TestList(t *testing.T) {
	names := List()

	assert.Equal(t, len(plc), len(names))
	for name, _ := range plc {
		assert.Contains(t, names, name)
	}
}

func TestWeight(t *testing.T) {
	cfg.SetServices(createServices())
	analytics := createAnalytics()
	cfg.SetNode(createNode())

	assert.InEpsilon(t, 0.0, plc["scalein"].Weight("service1", analytics), c_EPSILON)
	assert.InEpsilon(t, 0.16, plc["scalein"].Weight("service2", analytics), c_EPSILON)
	assert.InEpsilon(t, 0.0, plc["scalein"].Weight("service3", analytics), c_EPSILON)

	assert.InEpsilon(t, 0.0, plc["scaleout"].Weight("service1", analytics), c_EPSILON)
	assert.InEpsilon(t, 0.0, plc["scaleout"].Weight("service2", analytics), c_EPSILON)
	assert.InEpsilon(t, 0.25, plc["scaleout"].Weight("service3", analytics), c_EPSILON)
	assert.InEpsilon(t, 0.0, plc["scaleout"].Weight("service4", analytics), c_EPSILON)
}

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

	services := []cfg.Service{srv1, srv2, srv3, srv4}

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
	srv3A.Resources.Available = 1.0

	srv4A := analyzer.ServiceAnalytics{}
	srv4A.Resources.Available = 0.0

	analytics.Service = map[string]analyzer.ServiceAnalytics{
		"service1": srv1A,
		"service2": srv2A,
		"service3": srv3A,
		"service4": srv4A,
	}

	return analytics
}

func createNode() cfg.Node {
	n := cfg.Node{}
	n.Constraints.BaseServices = []string{"service3"}
	return n
}
