package strategy

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/data"
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
	var plc *data.Policy
	targets := []string{"pippo"}
	actions := map[string][]enum.Action{
		"pippo": []enum.Action{enum.START},
	}

	policies := []data.Policy{
		data.CreateMockPolicy("p", 0.0, targets, actions),
		data.CreateMockPolicy("p", 0.2, targets, actions),
		data.CreateMockPolicy("p", 0.5, targets, actions),
		data.CreateMockPolicy("p", 0.8, targets, actions),
		data.CreateMockPolicy("p", 1.0, targets, actions),
	}

	New("dummy")
	plc = MakeDecision(policies)
	assert.Equal(t, plc.Weight, 1.0)

	policies = []data.Policy{
		data.CreateMockPolicy("p", 0.0, targets, actions),
		data.CreateMockPolicy("p", 0.0, targets, actions),
		data.CreateMockPolicy("p", 0.0, targets, actions),
		data.CreateMockPolicy("p", 0.2, targets, actions),
	}

	New("probabilistic")
	plc = MakeDecision(policies)
	assert.Equal(t, plc.Weight, 0.2)
}
