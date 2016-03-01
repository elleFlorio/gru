package utils

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func TestGetCompleteRange(t *testing.T) {
	var err error
	var result []string

	invalid := "pippo"
	invalidLimit := "pippo-10"
	onveValue := "5"
	valid := "0-10"

	_, err = GetCompleteRange(invalid)
	assert.Error(t, err)

	_, err = GetCompleteRange(invalidLimit)
	assert.Error(t, err)

	result, err = GetCompleteRange(onveValue)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Contains(t, result, "5")

	result, err = GetCompleteRange(valid)
	assert.NoError(t, err)
	assert.Len(t, result, 11)
	assert.Contains(t, result, "0")
	assert.Contains(t, result, "10")
	assert.NotContains(t, result, "11")

}
