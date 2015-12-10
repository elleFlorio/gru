package agent

import (
	"errors"
	"strconv"
	"strings"
)

const c_CONFIG_PATH = "config/"

var ErrNoField error = errors.New("Unrecognized agent configuration field")

func parseConfig(agentData map[string]string) (GruAgentConfig, error) {
	dockerData := map[string]string{}
	autonomicData := map[string]string{}
	storageData := map[string]string{}
	metricData := map[string]string{}
	metricConfigData := map[string]interface{}{}
	for k, v := range agentData {
		data := k[strings.LastIndex(k, c_CONFIG_PATH):]
		fields := strings.Split(data, "/")
		switch fields[1] {
		case "docker":
			dockerData[fields[2]] = v
		case "autonomic":
			autonomicData[fields[2]] = v
		case "storage":
			storageData[fields[2]] = v
		case "metric":
			if fields[2] == "configuration" {
				metricConfigData[fields[3]] = v
			} else {
				metricData[fields[2]] = v
			}
		default:
			return GruAgentConfig{}, ErrNoField
		}
	}

	agentConfig := GruAgentConfig{
		Docker:    createDockerStruct(dockerData),
		Autonomic: createAutonomicStruct(autonomicData),
		Storage:   createStorageStruct(storageData),
		Metric:    createMetricStruct(metricData, metricConfigData),
	}

	return agentConfig, nil
}

func createDockerStruct(dockerData map[string]string) DockerConfig {
	docker := DockerConfig{}
	docker.DaemonUrl = dockerData["daemonurl"]
	docker.DaemonTimeout, _ = strconv.Atoi(dockerData["daemontimeout"])

	return docker
}

func createAutonomicStruct(autonomicData map[string]string) AutonomicConfig {
	autonomic := AutonomicConfig{}
	autonomic.LoopTimeInterval, _ = strconv.Atoi(autonomicData["looptimeinterval"])
	autonomic.MaxFriends, _ = strconv.Atoi(autonomicData["maxfriends"])

	return autonomic
}

func createStorageStruct(storageData map[string]string) StorageConfig {
	storage := StorageConfig{}
	storage.StorageService = storageData["storageservice"]

	return storage
}

func createMetricStruct(metricData map[string]string, metricConfigData map[string]interface{}) MetricConfig {
	metric := MetricConfig{}
	metric.MetricService = metricData["metricservice"]
	metric.Configuration = metricConfigData

	return metric
}
