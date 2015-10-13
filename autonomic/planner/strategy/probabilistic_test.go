package strategy

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/service"
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

func TestWeightedRandomElement(t *testing.T) {
	var plan *GruPlan
	var err error
	plans_empty := []GruPlan{}
	plans_zero := []GruPlan{CreateMockPlan(enum.WHITE, service.Service{}, []enum.Action{})}
	plans := CreateRandomPlans(10)

	plan, err = weightedRandomElement(plans_empty)
	assert.Error(t, err)

	plan, err = weightedRandomElement(plans_zero)
	assert.Error(t, err)

	plan, err = weightedRandomElement(plans)
	assert.NoError(t, err)
	assert.NotContains(t, plan.Actions, enum.NOACTION)

}
