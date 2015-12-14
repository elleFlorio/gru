package agent

import (
	"github.com/elleFlorio/gru/autonomic"
)

var agent GruAgentConfig

func Initialize(agentConfig GruAgentConfig) {
	agent = agentConfig
}

func GetAgent() GruAgentConfig {
	return agent
}

func Run() {
	startAutonomicManager()
}

func startAutonomicManager() {
	autonomic.Initialize(
		agent.Autonomic.LoopTimeInterval,
		agent.Autonomic.MaxFriends)
	autonomic.RunLoop()
}
