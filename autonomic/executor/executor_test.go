package executor

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/autonomic/planner/strategy"
	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/service"
	"github.com/elleFlorio/gru/storage"
)

func init() {
	storage.New("internal")
}

func TestRetrievePlan(t *testing.T) {
	defer storage.DeleteAllData(enum.PLANS)
	var err error

	_, err = retrievePlan()
	assert.Error(t, err)

	strategy.StoreMockPlan(enum.ORANGE, service.Service{}, []enum.Action{})
	_, err = retrievePlan()
	assert.NoError(t, err)
}

func TestRun(t *testing.T) {
	defer storage.DeleteAllData(enum.PLANS)
	srv := service.Service{
		Name:          "n",
		Image:         "i",
		Configuration: createServConfig(),
	}

	assert.NotPanics(t, Run)
	strategy.StoreMockPlan(enum.RED, srv, []enum.Action{enum.START})
	assert.Panics(t, Run)
	strategy.StoreMockPlan(enum.RED, srv, []enum.Action{enum.NOACTION})
	assert.NotPanics(t, Run)
}

func createServConfig() service.Config {
	cfg := service.Config{}
	cfg.Cmd = []string{"a", "b"}
	cfg.CpuSet = 4
	cfg.CpuShares = 512
	cfg.Entrypoint = []string{"d", "e"}
	cfg.Links = []string{"pippo"}
	cfg.Memory = "1G"
	cfg.PortBindings = createMockBindings()

	return cfg
}

func createMockBindings() map[string][]service.PortBinding {
	bindings := make(map[string][]service.PortBinding)
	bind_1 := service.PortBinding{"a", "b"}
	bind_2 := service.PortBinding{"c", "d"}
	bind_3 := service.PortBinding{"e", "f"}

	bindings["a"] = []service.PortBinding{bind_1, bind_2}
	bindings["b"] = []service.PortBinding{bind_3}

	return bindings
}
