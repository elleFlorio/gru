package enum

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func TestGetLabel(t *testing.T) {
	value_w := 0.2
	value_g := 0.4
	value_y := 0.6
	value_o := 0.8
	value_r := 1.0

	label := GetLabel(value_w)
	assert.Equal(t, WHITE, label, "(0.2) label should be white")

	label = GetLabel(value_g)
	assert.Equal(t, GREEN, label, "(0.4) label should be green")

	label = GetLabel(value_y)
	assert.Equal(t, YELLOW, label, "(0.6) label should be yellow")

	label = GetLabel(value_o)
	assert.Equal(t, ORANGE, label, "(0.8) label should be orange")

	label = GetLabel(value_r)
	assert.Equal(t, RED, label, "(1.0) label should be red")
}
