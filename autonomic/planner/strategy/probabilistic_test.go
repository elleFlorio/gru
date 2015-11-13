package strategy

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func TestRandomUniform(t *testing.T) {
	val := randUniform(0, 1)

	assert.InDelta(t, 0.5, val, 0.5)
}

func TestShuffle(t *testing.T) {
	plans_s := CreateRandomPlans(10)
	plans := make([]GruPlan, len(plans_s), len(plans_s))
	copy(plans, plans_s)

	shuffle(plans_s)
	assert.NotEqual(t, plans, plans_s)
}
