package agent

import (
	"github.com/elleFlorio/gru/autonomic"
)

var config GruAgentConfig

func Initialize(agentConfig GruAgentConfig) {
	config = agentConfig
}

func Config() GruAgentConfig {
	return config
}

func Run() {
	startAutonomicManager()
}

func startAutonomicManager() {
	autonomic.Initialize(
		config.Autonomic.LoopTimeInterval,
		config.Autonomic.MaxFriends)
	autonomic.RunLoop()
}
