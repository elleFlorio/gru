package agent

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func TestLoadGruAgentConfig(t *testing.T) {
	tmpfile := createMockConfigFile()
	defer os.Remove(tmpfile)

	err := LoadGruAgentConfig(tmpfile)

	if assert.NoError(t, err, "Loading should be done without errors") {
		assert.Equal(t, "unix:///var/run/docker.sock", config.Docker.DaemonUrl, "Expected different daemon url")
		assert.Equal(t, 10, config.Docker.DaemonTimeout, "Daemon timuout should be 0")
		assert.Equal(t, 5, config.Autonomic.LoopTimeInterval, "Loop time interval should be 5")
		assert.Equal(t, "/gru/config/services", config.Service.ServiceConfigFolder, "Expected different service config folder")
	}

}

func createMockConfigFile() string {

	mockConfigFile := `{
		"Service": {
			"ServiceConfigFolder":"/gru/config/services"
		},

		"Node": {
			"NodeConfigFile":"/gru/config/nodeconfig.json"
		},

		"Network": {
			"IpAddress":"127.0.0.1",
			"Port":"5000"
		},

		"Docker": {
			"DaemonUrl":"unix:///var/run/docker.sock",
			"DaemonTimeout":10
		},

		"Autonomic": {
			"LoopTimeInterval":5,
			"MaxFrineds":5,
			"DataToShare":"stats"
		},

		"Discovery": {
			"DiscoveryService":"etcd",
			"DiscoveryServiceUri":"http://127.0.0.1:4001"
		},
		
		"Storage": {
			"StorageService":"internal"
		}
	}`

	tmpfile, err := ioutil.TempFile(".", "gru_test_agent_config")
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile(tmpfile.Name(), []byte(mockConfigFile), 0644)

	return tmpfile.Name()

}
