package planner

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func TestSetPlannerStrategy(t *testing.T) {
	probabilistic := "probabilistic"
	notSupported := "notsupported"

	SetPlannerStrategy(probabilistic)
	assert.Equal(t, probabilistic, currentStrategy.Name(), "(probabilistic) Current strategy should be probabilistic")

	SetPlannerStrategy(notSupported)
	assert.Equal(t, "dummy", currentStrategy.Name(), "(notsupported) Current strategy should be dummy")
}
