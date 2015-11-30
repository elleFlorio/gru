package action

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/enum"
)

func TestGet(t *testing.T) {
	assert.Equal(t, &NoAction{}, Get(enum.NOACTION))
	assert.Equal(t, &Start{}, Get(enum.START))
	assert.Equal(t, &Stop{}, Get(enum.STOP))
}

func TestList(t *testing.T) {
	actionsList := List()

	assert.Len(t, actionsList, len(actions))
}

func TestRun(t *testing.T) {
	var err error

	noAction := &NoAction{}
	//start := &Start{}
	stop := &Stop{}
	config := GruActionConfig{}

	err = noAction.Run(config)
	assert.NoError(t, err)
	err = stop.Run(config)
	assert.Error(t, err)

	// Need to create a mock
	// container to complete the test
}
