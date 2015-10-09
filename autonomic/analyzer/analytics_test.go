package analyzer

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func TestConvertAnalyticsToData(t *testing.T) {
	analytics_ok := CreateMockAnalytics()

	_, err := convertAnalyticsToData(analytics_ok)
	assert.NoError(t, err, "(ok) stats convertion should produce no error")
}

func TestConvertDataToAnalytics(t *testing.T) {
	data_ok, err := convertAnalyticsToData(CreateMockAnalytics())
	data_bad := []byte{}

	_, err = ConvertDataToAnalytics(data_ok)
	assert.NoError(t, err, "(ok) data convertion should produce no error")

	_, err = ConvertDataToAnalytics(data_bad)
	assert.Error(t, err, "(bad) data convertion should produce an error")
}
