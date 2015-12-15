package action

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/service"
)

func createServConfig() service.Config {
	cfg := service.Config{}
	cfg.Cmd = []string{"a", "b"}
	cfg.CpusetCpus = "0"
	cfg.CpuShares = 512
	cfg.Entrypoint = []string{"d", "e"}
	cfg.Links = []string{"pippo"}
	cfg.Memory = "1G"
	cfg.PortBindings = createMockBindings()
	cfg.StopTimeout = 30

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

func TestCreateHostConfig(t *testing.T) {
	cfg := createServConfig()
	hostCfg := CreateHostConfig(cfg)

	assert.NotEmpty(t, hostCfg)
}

func TestCreatePortBindings(t *testing.T) {
	cfg := createServConfig()
	binding := createPortBindings(cfg.PortBindings)

	assert.Len(t, binding, len(cfg.PortBindings))
	for key, value := range cfg.PortBindings {
		assert.Len(t, binding[key], len(value))
	}
}

func TestCreateContainerConfig(t *testing.T) {
	cfg := createServConfig()
	contCfg := CreateContainerConfig(cfg)

	assert.NotEmpty(t, contCfg)
}
