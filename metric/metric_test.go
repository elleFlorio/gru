package metric

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/autonomic/analyzer"
	"github.com/elleFlorio/gru/autonomic/monitor"
	"github.com/elleFlorio/gru/autonomic/planner/strategy"
	cfg "github.com/elleFlorio/gru/configuration"
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

	UpdateMetrics()
	assert.Equal(t, 0.0, Metrics().Service["service1"].Stats.CpuTot)
	assert.Equal(t, 0.0, Metrics().Service["service2"].Analytics.Cpu)
	assert.Equal(t, "noaction", Metrics().Plan.Policy)

	monitor.StoreMockStats()
	analyzer.StoreMockAnalytics()
	mockService, _ := service.GetServiceByName("service1")
	strategy.StoreMockPlan("policy", 1.0, *mockService, enum.Actions{enum.START})
	UpdateMetrics()
	assert.Equal(t, 0.7, Metrics().Service["service1"].Stats.CpuTot)
	assert.Equal(t, 0.2, Metrics().Service["service2"].Analytics.Cpu)
	assert.Equal(t, "policy", Metrics().Plan.Policy)
}

func TestCreateInfluxMetrics(t *testing.T) {
	New("influxdb", createInfluxMockConfig())
	mockMetrics := CreateMockMetrics()
	points, err := createInfluxMetrics(mockMetrics)
	assert.NoError(t, err)
	assert.NotEmpty(t, points)
}
