package logreader

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func TestLogReading(t *testing.T) {
	logFile := createMockLog()
	logStream, _ := os.Open(logFile)
	defer os.Remove(logFile)
	services := []string{"service1", "service2"}

	Initialize(services)
	StartLogReader()
	StartCollector(logStream)
	time.Sleep(1 * time.Second)
	metricsSrv1 := srvMetrics["service1"].userDef
	metricsSrv2 := srvMetrics["service2"].userDef
	assert.Len(t, metricsSrv1, 2)
	assert.Len(t, metricsSrv2, 1)
	m1Srv1 := metricsSrv1["metric1"]
	m2Srv1 := metricsSrv1["metric2"]
	m1Srv2 := metricsSrv2["metric1"]
	m2Srv2 := metricsSrv2["metric2"]
	assert.Len(t, m1Srv1.GetValues(), 2)
	assert.Len(t, m2Srv1.GetValues(), 1)
	assert.Len(t, m1Srv2.GetValues(), 1)
	assert.Nil(t, m2Srv2.GetValues())
	Stop()
}

func createMockLog() string {
	noData := "no data\n"
	gruService1Metric1 := "gru:service1:metric1:1:ms\n"
	gruService1Metric2 := "gru:service1:metric2:2:ms\n"
	gruService2Metric1 := "gru:service2:metric1:3:ms\n"
	gruWrongFloat := "gru:service2:metric2:1,4:ms\n"
	gruWrongFormat := "gru:pippo:1\n"

	mockLog := noData +
		gruService1Metric1 +
		gruService1Metric1 +
		gruWrongFloat +
		gruService1Metric2 +
		gruWrongFormat +
		gruService2Metric1

	tmpfile, err := ioutil.TempFile(".", "gru_test_log_parser")
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile(tmpfile.Name(), []byte(mockLog), 0644)

	return tmpfile.Name()
}
