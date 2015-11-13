package policy

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/autonomic/analyzer"
	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/service"
)

var plc map[string]GruPolicy

func init() {
	plc = map[string]GruPolicy{
		"scalein":  &ScaleIn{},
		"scaleout": &ScaleOut{},
	}
}

func TestGetPolicies(t *testing.T) {
	pls := GetPolicies()
	names := make([]string, 0)
	actions := make([][]enum.Action, 0)
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

func TestLabel(t *testing.T) {
	services := createServices()
	service.UpdateServices(services)
	analytics := createAnalytics()

	assert.Equal(t, enum.GREEN, plc["scalein"].Label("service1", analytics))
	assert.Equal(t, enum.WHITE, plc["scalein"].Label("service2", analytics))

	assert.Equal(t, enum.WHITE, plc["scaleout"].Label("service1", analytics))
	assert.Equal(t, enum.GREEN, plc["scaleout"].Label("service2", analytics))
	assert.Equal(t, enum.WHITE, plc["scaleout"].Label("service3", analytics))
}

func createServices() []service.Service {
	srv1 := service.Service{}
	srv1.Name = "service1"
	srv1.Constraints.MaxActive = 1
	srv1.Instances.Running = []string{"instance1_1"}

	srv2 := service.Service{}
	srv2.Name = "service2"
	srv2.Constraints.MinActive = 2
	srv2.Instances.Running = []string{"instance2_1"}
	srv2.Instances.Pending = []string{"instance2_2"}

	srv3 := service.Service{}
	srv3.Name = "service3"
	srv3.Instances.Running = []string{"instance3_1"}

	services := []service.Service{srv1, srv2, srv3}

	return services
}

func createAnalytics() analyzer.GruAnalytics {
	analytics := analyzer.GruAnalytics{}

	srv1A := analyzer.ServiceAnalytics{}
	srv1A.Load = enum.RED
	srv1A.Resources.Cpu = enum.YELLOW

	srv2A := analyzer.ServiceAnalytics{}
	srv2A.Load = enum.YELLOW
	srv2A.Resources.Cpu = enum.GREEN

	srv3A := analyzer.ServiceAnalytics{}
	srv3A.Load = enum.GREEN
	srv3A.Resources.Cpu = enum.GREEN
	srv3A.Resources.Available = enum.RED

	analytics.Service = map[string]analyzer.ServiceAnalytics{
		"service1": srv1A,
		"service2": srv2A,
		"service3": srv3A,
	}

	return analytics
}
