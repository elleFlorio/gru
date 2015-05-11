package cli

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadGruAgentConfig(t *testing.T) {
	tmpfile := createMockConfigFile()
	defer os.Remove(tmpfile)

	config, err := LoadGruAgentConfig(tmpfile)

	if assert.NoError(t, err, "Loading should be done without errors") {
		assert.Equal(t, "unix://var/run/docker.sock", config.DaemonUrl, "Expected different daemon url")
		assert.Equal(t, 0, config.DaemonTimeout, "Daemon timuout should be 0")
		assert.Equal(t, 5, config.LoopTimeInterval, "Loop time interval should be 5")
		assert.Equal(t, "/config/service", config.ServiceConfigFolder, "Expected different service config folder")
	}

}

func createMockConfigFile() string {

	mockConfigFile := `{
        "DaemonUrl":"unix://var/run/docker.sock",
        "DaemonTimeout":0,
        "LoopTimeInterval":5,
        "ServiceConfigFolder":"/config/service"
    }`

	tmpfile, err := ioutil.TempFile(".", "gru_test_agent_config")
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile(tmpfile.Name(), []byte(mockConfigFile), 0644)

	return tmpfile.Name()

}
