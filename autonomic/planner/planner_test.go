package planner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPlanner(t *testing.T) {
	c_err := make(chan error)

	planDummy := NewPlanner("dummy", c_err)
	assert.Equal(t, "dummy", planDummy.strtg.Name(), "Strategy should be dummy")

	planError := NewPlanner("notImplementedStrategy", c_err)
	assert.Equal(t, "dummy", planError.strtg.Name(), "Default strategy in case of error should be dummy")

}
