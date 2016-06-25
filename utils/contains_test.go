package utils

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func TestContainsString(t *testing.T) {
	values := []string{"a", "b", "c"}
	valueTrue := "a"
	valueFalse := "d"

	assert.True(t, ContainsString(values, valueTrue))
	assert.False(t, ContainsString(values, valueFalse))
}
