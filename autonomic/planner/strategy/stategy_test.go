package strategy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	strategy, err := New("dummy")
	assert.NoError(t, err, "Dummy strategy should be created without errors")
	assert.Equal(t, "dummy", strategy.Name(), "The name of the strategy should be 'dummy'")

	strategy, err = New("notImplemented")
	assert.Error(t, err, "Strategies not implemented should raise and error")
}

func TestList(t *testing.T) {
	names := List()
	assert.Contains(t, names, "dummy", "The list of strategies should contain 'dummy'")
}
