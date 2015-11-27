package executor

import (
	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/action"
	"github.com/elleFlorio/gru/autonomic/planner/strategy"
	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/service"
)

func Run(plan *strategy.GruPlan) {
	log.Debugln("Running Executor")

	if plan == nil {
		log.WithField("error", "No plan to execute").Errorln("Cannot execute actions.")
	} else {
		config := buildConfig(plan.Target)
		executeActions(plan.Actions, config)
	}
}

func buildConfig(srv *service.Service) action.GruActionConfig {
	actConfig := action.GruActionConfig{}
	actConfig.Service = srv.Name
	actConfig.Instances = srv.Instances
	actConfig.HostConfig = action.CreateHostConfig(srv.Configuration)
	actConfig.ContainerConfig = action.CreateContainerConfig(srv.Configuration)
	actConfig.ContainerConfig.Image = srv.Image
	actConfig.Parameters.StopTimeout = srv.Configuration.StopTimeout

	return actConfig
}

func executeActions(actions []enum.Action, config action.GruActionConfig) {
	var err error
	for _, actionType := range actions {
		act := action.Get(actionType)
		err = act.Run(config)
		if err != nil {
			log.WithFields(log.Fields{
				"error":  err,
				"action": act.Type().ToString(),
			}).Errorln("Action not executed")
		}

		log.WithFields(log.Fields{
			"target": config.Service,
			"action": act.Type().ToString(),
		}).Infoln("Action executed")
	}
}
