package strategy

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/service"
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
	var plan *GruPlan

	s := service.Service{}
	a := []enum.Action{enum.START}

	plans := []GruPlan{
		CreateMockPlan(0.0, s, a),
		CreateMockPlan(0.2, s, a),
		CreateMockPlan(0.5, s, a),
		CreateMockPlan(0.8, s, a),
		CreateMockPlan(1.0, s, a),
	}
	New("dummy")
	plan = MakeDecision(plans)
	assert.Equal(t, plan.Weight, 1.0)

	plans = []GruPlan{
		CreateMockPlan(0.0, s, a),
		CreateMockPlan(0.0, s, a),
		CreateMockPlan(0.0, s, a),
		CreateMockPlan(0.2, s, []enum.Action{enum.NOACTION}),
	}
	New("probabilistic")
	plan = MakeDecision(plans)
	assert.Equal(t, plan.Weight, 0.2)

}
