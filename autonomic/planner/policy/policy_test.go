package policy

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	//"github.com/elleFlorio/gru/autonomic/analyzer"
	"github.com/elleFlorio/gru/enum"
)

var plc map[string]GruPolicy

func init() {
	plc = map[string]GruPolicy{
		"scalein":  &ScaleIn{},
		"scaleout": &ScaleOut{},
	}
}

func TestGetPolicies(t *testing.T) {
	pls := GetPolicies()
	names := make([]string, 0)
	actions := make([][]enum.Action, 0)
	for _, item := range pls {
		names = append(names, item.Name())
		actions = append(actions, item.Actions())
	}

	assert.Equal(t, len(plc), len(pls))
	for _, item := range plc {
		assert.Contains(t, names, item.Name())
		assert.Contains(t, actions, item.Actions())
	}
}

func TestList(t *testing.T) {
	names := List()

	assert.Equal(t, len(plc), len(names))
	for name, _ := range plc {
		assert.Contains(t, names, name)
	}
}

func TestLabel(t *testing.T) {
	//TODO
}
