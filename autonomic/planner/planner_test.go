package planner

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func TestSetPlannerStrategy(t *testing.T) {
	probcumulative := "probcumulative"
	probdelta := "probdelta"
	notSupported := "notsupported"

	SetPlannerStrategy(probcumulative)
	assert.Equal(t, probcumulative, currentStrategy.Name(), "(probcumulative) Current strategy should be probcumulative")

	SetPlannerStrategy(probdelta)
	assert.Equal(t, probdelta, currentStrategy.Name(), "(probdelta) Current strategy should be probdelta")

	SetPlannerStrategy(notSupported)
	assert.Equal(t, "dummy", currentStrategy.Name(), "(notsupported) Current strategy should be dummy")
}
