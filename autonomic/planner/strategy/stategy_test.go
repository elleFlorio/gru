package strategy

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/autonomic/planner/policy"
	"github.com/elleFlorio/gru/enum"
)

func TestNew(t *testing.T) {
	var err error

	_, err = New("dummy")
	assert.NoError(t, err)
	assert.Equal(t, "dummy", Name())

	_, err = New("probabilistic")
	assert.NoError(t, err)
	assert.Equal(t, "probabilistic", Name())

	_, err = New("notImplemented")
	assert.Error(t, err)
	assert.Equal(t, "dummy", Name())
}

func TestList(t *testing.T) {
	names := List()
	assert.Contains(t, names, "dummy")
	assert.Contains(t, names, "probabilistic")
}

func TestMakeDecision(t *testing.T) {
	var plc *policy.Policy

	targets := map[string][]enum.Action{
		"pippo": []enum.Action{enum.START},
	}

	policies := []policy.Policy{
		policy.CreateMockPolicy("p", 0.0, targets),
		policy.CreateMockPolicy("p", 0.2, targets),
		policy.CreateMockPolicy("p", 0.5, targets),
		policy.CreateMockPolicy("p", 0.8, targets),
		policy.CreateMockPolicy("p", 1.0, targets),
	}

	New("dummy")
	plc = MakeDecision(policies)
	assert.Equal(t, plc.Weight, 1.0)

	policies = []policy.Policy{
		policy.CreateMockPolicy("p", 0.0, targets),
		policy.CreateMockPolicy("p", 0.0, targets),
		policy.CreateMockPolicy("p", 0.0, targets),
		policy.CreateMockPolicy("p", 0.2, targets),
	}

	New("probabilistic")
	plc = MakeDecision(policies)
	assert.Equal(t, plc.Weight, 0.2)
}
