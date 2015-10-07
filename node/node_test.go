package node

import (
	"os"
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func TestLoadNodeConfig(t *testing.T) {
	tmpFile := createMockConfigFileNode()

	err := LoadNodeConfig(tmpFile)
	config := Config()

	if assert.NoError(t, err, "Node config loading should be done without errors") {
		assert.NotEmpty(t, config.UUID, "Node UUID should be set")
		assert.Equal(t, "mockNode", config.Name, "Expected different node name")
		assert.Equal(t, 0.8, config.Constraints.CpuMax, "Node cpu max should be 0.8")
		assert.Equal(t, 0.2, config.Constraints.CpuMin, "Node cpu min should be 0.2")
		assert.Equal(t, 10, config.Constraints.MaxInstances, "Node max instances should be 10")
	}

	UpdateNodeConfig(Node{})
	os.Remove(tmpFile)
}
