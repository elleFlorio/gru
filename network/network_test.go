package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPort(t *testing.T) {
	_, err := GetPort()
	assert.NoError(t, err, "Port retrieval should generate no error")
}

func TestGetHostIp(t *testing.T) {
	_, err := GetHostIp()
	assert.NoError(t, err, "IP retrieval should generate no error")
}
