package agent

import (
	"github.com/elleFlorio/gru/autonomic"
	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/data"
)

func Initialize() {
	autonomic.UpdatePlannerStrategy(cfg.GetAgentAutonomic().PlannerStrategy)
	data.InitializeFriendsData(cfg.GetAgentCommunication().MaxFriends)
}

func StartMonitoring() {
	autonomic.Start()
}

func Run() {
	startAutonomicManager()
}

func UpdateStrategy() {
	autonomic.UpdatePlannerStrategy(cfg.GetAgentAutonomic().PlannerStrategy)
}

func startAutonomicManager() {
	autonomic.RunLoop()
}
