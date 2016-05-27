package strategy

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/data"
	"github.com/elleFlorio/gru/enum"
)

func TestWeightedRandomElement(t *testing.T) {
	var plc *data.Policy
	var threshold float64
	var policies []data.Policy

	targets := []string{"pippo"}
	actions := map[string][]enum.Action{
		"pippo": []enum.Action{enum.START},
	}

	policies = []data.Policy{
		data.CreateMockPolicy("p1", 0.2, targets, actions),
		data.CreateMockPolicy("p2", 0.4, targets, actions),
		data.CreateMockPolicy("p3", 0.6, targets, actions),
		data.CreateMockPolicy("p4", 0.2, targets, actions),
	}

	threshold = 0.4
	plc = weightedRandomElement(policies, threshold)
	assert.Equal(t, "p2", plc.Name)

	threshold = 0.8
	plc = weightedRandomElement(policies, threshold)
	assert.Equal(t, "p3", plc.Name)
}
