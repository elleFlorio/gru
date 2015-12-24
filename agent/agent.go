package agent

import (
	"github.com/elleFlorio/gru/autonomic"
	cfg "github.com/elleFlorio/gru/configuration"
)

func Initialize() {
	autonomic.Initialize(
		cfg.GetAgentAutonomic().PlannerStrategy,
	)
}

func StartMonitoring() {
	autonomic.Start()
}

func Run() {
	startAutonomicManager()
}

func startAutonomicManager() {
	autonomic.RunLoop(
		cfg.GetAgentAutonomic().LoopTimeInterval,
	)
}
