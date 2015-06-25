package executor

import (
	"testing"

	"github.com/samalba/dockerclient"
	"github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/autonomic/planner/strategy"
	"github.com/elleFlorio/gru/service"
)

var mockExec executor

func init() {
	mockExec = executor{}
}

func TestBuildActions(t *testing.T) {
	mockPlans := strategy.CreateMockPlans(0.8, 0.3, 0.5)
	planStart := mockPlans[0]
	planErr := mockPlans[1]

	// Correct
	actions, err := mockExec.buildActions(&planStart)
	assert.NoError(t, err, "planStart should produce no errors")
	assert.Equal(t, 2, len(actions), "Created actions should have length 1")
	assert.Equal(t, "start", actions[0].Name(), "Created action should have name 'start'")

	// Error
	actions, err = mockExec.buildActions(&planErr)
	assert.Error(t, err, "planErr should produce an error")
	assert.Equal(t, 0, len(actions), "Created actions should have length 0")

}

func TestBuildConfig(t *testing.T) {
	mockPlans := strategy.CreateMockPlans(0.8, 0.3, 0.5)
	planStart := mockPlans[0]
	srv := service.Service{}
	mockDocker, _ := dockerclient.NewDockerClient("daemon_url", nil)

	config := mockExec.buildConfig(&planStart, &srv, mockDocker)
	assert.Equal(t, "service1", config.Service, "Configuration service should be 'service1'")
}
