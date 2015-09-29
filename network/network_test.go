package network

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func TestGetPort(t *testing.T) {
	_, err := getPort()
	assert.NoError(t, err, "Port retrieval should generate no error")
}

func TestGetHostIp(t *testing.T) {
	_, err := getHostIp()
	assert.NoError(t, err, "IP retrieval should generate no error")
}
