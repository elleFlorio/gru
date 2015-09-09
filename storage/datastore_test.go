package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	supported := "internal"
	notSuported := "notSupported"

	test, err := New(supported)
	assert.NoError(t, err, "Supported storage should produce no error")
	assert.Equal(t, "internal", test.Name(), "Storage should be 'internal'")

	test = DataStore()
	assert.Equal(t, "internal", test.Name(), "(supported) Retrieved datastore should be 'internal'")

	test, err = New(notSuported)
	assert.Error(t, err, "Not supported storage should produce an error")
	assert.Equal(t, "internal", test.Name(), "If storage is not supported the default one should be 'internal'")

	test = DataStore()
	assert.Equal(t, "internal", test.Name(), "(not supported) retrieved datastore should be 'internal'")
}
