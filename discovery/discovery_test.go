package discovery

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	correct := "etcd"
	notSupported := "notSupported"

	dscvr, err := New(correct, "http://localhost:5000")
	assert.NoError(t, err, "etcd should be supported")
	assert.Equal(t, "etcd", dscvr.Name(), "Discovery name should be 'etcd'")

	dscvr, err = New(notSupported, "http://localhost:5000")
	assert.Error(t, err, "Not supported service should produce an error")
	assert.Equal(t, "noservice", dscvr.Name(), "(not supported) Discovery name should be 'noservice'")

}
