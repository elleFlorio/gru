package metric

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func TestMetricManager(t *testing.T) {
	logFile := createMockLog()
	logStream, _ := os.Open(logFile)
	defer os.Remove(logFile)

	Manager().Start()
	Manager().StartCollector(logStream)
	time.Sleep(1 * time.Second)
	data := Manager().GetMetrics()
	assert.Len(t, data, 2)
	Manager().Stop()
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
