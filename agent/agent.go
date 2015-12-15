package agent

import (
	"github.com/elleFlorio/gru/autonomic"
	cfg "github.com/elleFlorio/gru/configuration"
)

func Run() {
	startAutonomicManager()
}

func startAutonomicManager() {
	autonomic.Initialize(
		cfg.GetAgentAutonomic().LoopTimeInterval,
		cfg.GetAgentAutonomic().MaxFriends,
	)
	autonomic.RunLoop()
}
