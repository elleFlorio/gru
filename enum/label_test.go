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

func TestValueFrom(t *testing.T) {
	expected_w := 0.0
	expected_g := 0.2
	expected_y := 0.4
	expected_o := 0.6
	expected_r := 0.8

	assert.Equal(t, expected_w, ValueFrom(WHITE))
	assert.Equal(t, expected_g, ValueFrom(GREEN))
	assert.Equal(t, expected_y, ValueFrom(YELLOW))
	assert.Equal(t, expected_o, ValueFrom(ORANGE))
	assert.Equal(t, expected_r, ValueFrom(RED))
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
