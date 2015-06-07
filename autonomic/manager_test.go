package autonomic

import (
	"testing"

	"github.com/samalba/dockerclient"
	"github.com/stretchr/testify/assert"
)

func TestNewAutoManager(t *testing.T) {
	mockDocker, _ := dockerclient.NewDockerClient("daemonUrl", nil)
	mockAutoManager := NewAutoManager(mockDocker, 5)

	assert.Equal(t, 5, mockAutoManager.LoopTimeInterval, "Time interval should be 5 secs")
}
