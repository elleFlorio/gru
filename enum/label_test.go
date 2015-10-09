package enum

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func TestFromValue(t *testing.T) {
	value_w := 0.2
	value_g := 0.4
	value_y := 0.6
	value_o := 0.8
	value_r := 1.0

	assert.Equal(t, WHITE, FromValue(value_w))
	assert.Equal(t, GREEN, FromValue(value_g))
	assert.Equal(t, YELLOW, FromValue(value_y))
	assert.Equal(t, ORANGE, FromValue(value_o))
	assert.Equal(t, RED, FromValue(value_r))
}

func TestFromLabelValue(t *testing.T) {
	l_value_w := -0.9
	l_value_g := -0.4
	l_value_y := 0.1
	l_value_o := 0.6
	l_value_r := 1.0

	assert.Equal(t, WHITE, FromLabelValue(l_value_w))
	assert.Equal(t, GREEN, FromLabelValue(l_value_g))
	assert.Equal(t, YELLOW, FromLabelValue(l_value_y))
	assert.Equal(t, ORANGE, FromLabelValue(l_value_o))
	assert.Equal(t, RED, FromLabelValue(l_value_r))

}
