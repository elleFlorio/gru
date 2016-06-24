package metric

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/data"
	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/service"
	"github.com/elleFlorio/gru/storage"
)

func init() {
	storage.New("internal")
}

func TestNew(t *testing.T) {
	noService := "noservice"
	noServiceConf := CreateMetricsMockConfig(noService)
	influxService := "influxdb"
	influxConf := CreateMetricsMockConfig(influxService)

	test, testErr := New(influxService, influxConf)
	assert.NoError(t, testErr)
	assert.Equal(t, test.Name(), influxService)
	influxConf["Username"] = 1
	test, testErr = New(influxService, influxConf)
	assert.Error(t, testErr)
	assert.Equal(t, test.Name(), noService)

	test, testErr = New(noService, noServiceConf)
	assert.Error(t, testErr)
	assert.Equal(t, test.Name(), noService)
}

func TestUpdateMetrics(t *testing.T) {
	cfg.SetServices(service.CreateMockServices())

	updateMetrics()
	assert.Equal(t, 0.0, Metrics().Service["service1"].Stats.BaseMetrics[enum.METRIC_CPU_AVG.ToString()])
	assert.Equal(t, 0.0, Metrics().Service["service2"].Analytics.BaseAnalytics[enum.METRIC_CPU_AVG.ToString()])
	assert.Equal(t, "noaction", Metrics().Policy.Name)

	data.SaveMockStats()
	data.SaveMockAnalytics()
	data.SaveSharedCluster(data.CreateMockShared())
	plc := data.CreateMockPolicy("policy", 1.0, []string{"pippo"}, map[string][]enum.Action{})
	data.SavePolicy(plc)
	updateMetrics()
	assert.Equal(t, 0.6, Metrics().Service["service1"].Stats.BaseMetrics[enum.METRIC_CPU_AVG.ToString()])
	assert.Equal(t, 0.1, Metrics().Service["service2"].Analytics.BaseAnalytics[enum.METRIC_CPU_AVG.ToString()])
	assert.Equal(t, 0.6, metrics.Service["service1"].Shared.BaseShared[enum.METRIC_CPU_AVG.ToString()])
	assert.Equal(t, "policy", Metrics().Policy.Name)
}

func TestCreateInfluxMetrics(t *testing.T) {
	New("influxdb", createInfluxMockConfig())
	mockMetrics := CreateMockMetrics()
	points, err := createInfluxMetrics(mockMetrics)
	assert.NoError(t, err)
	assert.NotEmpty(t, points)
}
