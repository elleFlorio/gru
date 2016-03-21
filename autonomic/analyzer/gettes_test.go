package analyzer

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/data"
	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/storage"
)

func init() {
	storage.New("internal")
}

func TestGetters(t *testing.T) {
	defer storage.DeleteAllData(enum.ANALYTICS)
	data.StoreMockAnalytics()
	analytics := data.CreateMockAnalytics()
	// data, _ := convertAnalyticsToData(analytics)
	// storage.StoreClusterData(data, enum.ANALYTICS)

	//assert.Equal(t, analytics, GetNodeAnalytics())
	assert.Equal(t, analytics.Service["service1"], GetServiceAnalytics("service1"))
	assert.Equal(t, analytics.Service, GetServicesAnalytics())
	assert.Equal(t, analytics.System, GetSystemAnalytics())
	assert.Equal(t, analytics.Cluster, GetClusterAnalytics())

}
