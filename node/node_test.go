package node

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadNodeConfig(t *testing.T) {
	// Reading Ip and port
	tmpFile := createMockConfigFileNode(true, true)

	err := LoadNodeConfig(tmpFile)
	config := Config()

	if assert.NoError(t, err, "Node config loading should be done without errors") {
		assert.NotEmpty(t, config.UUID, "Node UUID should be set")
		assert.Equal(t, "mockNode", config.Name, "Expected different node name")
		assert.Equal(t, "127.0.0.1", config.IpAddr, "Node Ip address should be set")
		assert.Equal(t, "8080", config.Port, "Node port should be set")
		assert.Equal(t, 0.8, config.Constraints.CpuMax, "Node cpu max should be 0.8")
		assert.Equal(t, 0.2, config.Constraints.CpuMin, "Node cpu min should be 0.2")
		assert.Equal(t, 10, config.Constraints.MaxInstances, "Node max instances should be 10")
	}

	UpdateNode(Node{})
	os.Remove(tmpFile)

	//Reading only Ip
	tmpFile = createMockConfigFileNode(true, false)

	err = LoadNodeConfig(tmpFile)
	config = Config()
	assert.NotEmpty(t, config.Port, "Port should be set by agent")

	UpdateNode(Node{})
	os.Remove(tmpFile)

	//Reading only port
	tmpFile = createMockConfigFileNode(false, true)

	err = LoadNodeConfig(tmpFile)
	config = Config()
	assert.NotEmpty(t, config.IpAddr, "Ip address should be set by agent")

	UpdateNode(Node{})
	os.Remove(tmpFile)
}
