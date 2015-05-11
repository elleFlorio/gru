package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const number int = 1

func TestList(t *testing.T) {
	actions := List()

	assert.Len(t, actions, number, "Number of current actions should be 2")

	assert.Equal(t, "start", actions[0], "The name of the first action should be 'start'")
}
