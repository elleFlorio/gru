package discovery

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	notSupported := "notSupported"

	dscvr, err := New(notSupported, "http://localhost:5000")
	assert.Error(t, err)
	assert.Equal(t, "noservice", dscvr.Name())
}
