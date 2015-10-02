package monitor

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func TestConvertStatsToData(t *testing.T) {
	stats_ok := CreateMockStats()

	_, err := convertStatsToData(stats_ok)
	assert.NoError(t, err, "(ok) stats convertion should produce no error")
}

func TestConvertDataToStats(t *testing.T) {
	data_ok, err := convertStatsToData(CreateMockStats())
	data_bad := []byte{}

	_, err = ConvertDataToStats(data_ok)
	assert.NoError(t, err, "(ok) data convertion should produce no error")

	_, err = ConvertDataToStats(data_bad)
	assert.Error(t, err, "(bad) data convertion should produce an error")
}
