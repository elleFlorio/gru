package executor

import (
	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/action"
	"github.com/elleFlorio/gru/autonomic/planner/strategy"
	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/service"
	"github.com/elleFlorio/gru/storage"
)

func Run() {
	log.WithField("status", "start").Infoln("Running Executor")
	defer log.WithField("status", "done").Infoln("Running Executor")

	plan, err := retrievePlan()
	if err != nil {
		log.WithField("error", "Cannot execute actions").Errorln("Running Executor.")
	} else {
		config := buildConfig(plan.Target)
		executeActions(plan.Actions, config)
	}
}

func retrievePlan() (strategy.GruPlan, error) {
	plan := strategy.GruPlan{}
	dataPlan, err := storage.GetLocalData(enum.PLANS)
	if err != nil {
		log.WithField("error", err).Errorln("Cannot retrieve plan data.")
	} else {
		plan, err = strategy.ConvertDataToPlan(dataPlan)
	}

	return plan, err
}

func buildConfig(srv *service.Service) action.GruActionConfig {
	actConfig := action.GruActionConfig{}
	actConfig.Service = srv.Name
	actConfig.Instances = srv.Instances
	actConfig.HostConfig = action.CreateHostConfig(srv.Configuration)
	actConfig.ContainerConfig = action.CreateContainerConfig(srv.Configuration)

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
