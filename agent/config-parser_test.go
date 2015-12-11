package agent

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {
	agentData := generateData()

	agentConf, err := parseConfig(agentData)
	assert.NoError(t, err)
	assert.Equal(t, 30, agentConf.Autonomic.LoopTimeInterval)
	assert.Equal(t, "influxdb", agentConf.Metric.MetricService)

	agentData["gru/pippo/config/field"] = "topolino"
	agentConf, err = parseConfig(agentData)
	assert.Error(t, err)

}

func generateData() map[string]string {
	configPath := "gru/pippo/config/"
	dockerPath1 := "docker/daemonurl"
	dockerPath2 := "docker/daemontimeout"
	autonomicPath1 := "autonomic/looptimeinterval"
	autonomicPath2 := "autonomic/maxfriends"
	storagePath1 := "storage/storageservice"
	metricPath1 := "metric/metricservice"
	metricConfPath1 := "metric/configuration/url"
	metricConfPath2 := "metric/configuration/something"

	data := map[string]string{
		configPath + dockerPath1:     "aaaaaaa",
		configPath + dockerPath2:     "10",
		configPath + autonomicPath1:  "30",
		configPath + autonomicPath2:  "5",
		configPath + storagePath1:    "internal",
		configPath + metricPath1:     "influxdb",
		configPath + metricConfPath1: "bbbbbbbbb",
		configPath + metricConfPath2: "ccccccccc",
	}

	return data

}
