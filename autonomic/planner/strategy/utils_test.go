package strategy

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/data"
)

func TestRandomUniform(t *testing.T) {
	val := randUniform(0, 1)

	assert.InDelta(t, 0.5, val, 0.5)
}

func TestShuffle(t *testing.T) {
	policies_s := data.CreateRandomMockPolicies(5)
	policies := make([]data.Policy, len(policies_s), len(policies_s))
	copy(policies, policies_s)

	shuffle(policies_s)
	assert.NotEqual(t, policies, policies_s)
}
