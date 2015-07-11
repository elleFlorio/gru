package node

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadNodeConfig(t *testing.T) {
	tmpFile := createMockConfigFileNode()
	defer os.Remove(tmpFile)

	err := LoadNodeConfig(tmpFile)
	config := GetNodeConfig()

	if assert.NoError(t, err, "Node config loading should be done without errors") {
		assert.Equal(t, "mockNode", config.Name, "Expected different node name")
		assert.Equal(t, 0.8, config.Constraints.CpuMax, "Node cpu max should be 0.8")
		assert.Equal(t, 0.2, config.Constraints.CpuMin, "Node cpu min should be 0.2")
		assert.Equal(t, 10, config.Constraints.MaxInstances, "Node max instances should be 10")
	}
}
