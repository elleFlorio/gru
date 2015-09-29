package utils

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

const uuidLen int = 32

func TestGenerateUUID(t *testing.T) {
	uuid, err := GenerateUUID()
	assert.NoError(t, err, "UUID generation should produce no errors")
	assert.Len(t, uuid, uuidLen, "generated UUID should be 32 characters")
}
