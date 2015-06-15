package action

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

const number int = 3

func TestNew(t *testing.T) {
	correct := "start"
	notImplemented := "notImplemented"

	act, _ := New(correct)
	assert.Equal(t, "start", act.Name(), "The name of the function should be 'start'")

	act, err := New(notImplemented)
	assert.Error(t, err, "If an action is not implemented an error should be raised")
}

func TestList(t *testing.T) {
	actions := List()

	assert.Len(t, actions, number, "Number of current actions should be "+strconv.Itoa(number))

	assert.Equal(t, "noAction", actions[0], "The name of the first action should be 'start'")
}
