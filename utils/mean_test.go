package utils

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func TestMean(t *testing.T) {
	values := []float64{2, 4, 6}
	valuesEmpty := []float64{}

	assert.Equal(t, 4.0, Mean(values))
	assert.Equal(t, 0.0, Mean(valuesEmpty))
}
