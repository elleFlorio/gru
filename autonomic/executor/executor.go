package executor

import (
	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/action"
	"github.com/elleFlorio/gru/autonomic/planner/strategy"
	"github.com/elleFlorio/gru/service"
)

//FIXME should build a configuration for each service
func Run(plan strategy.GruPlan) {
	log.WithField("status", "start").Debugln("Running executor")
	defer log.WithField("status", "done").Debugln("Running executor")

	config := &action.GruActionConfig{}
	actions, err := buildActions(plan)
	if err != nil {
		log.WithFields(log.Fields{
			"status": "error",
			"error":  err,
		}).Debugln("Running executor")
	}

	srv, err := service.GetServiceByName(plan.Service)

	if err == nil {
		config = buildConfig(plan, srv)
	}

	log.WithFields(log.Fields{
		"status": "config builded",
		"config": config,
	}).Debugln("Running executor")

	log.WithFields(log.Fields{
		"status":  "actuation",
		"actions": len(actions),
	}).Debugln("Running executor")

	for _, act := range actions {
		act.Run(config)
	}
}

func buildActions(plan strategy.GruPlan) ([]action.GruAction, error) {
	var err error
	actions := make([]action.GruAction, 0, len(plan.Actions))
	for _, name := range plan.Actions {
		act, err := action.New(name)
		if err != nil {
			return actions, err
		} else {
			actions = append(actions, act)
		}
	}

	return actions, err
}

// TODO container config should be configured properly
func buildConfig(plan strategy.GruPlan, srv *service.Service) *action.GruActionConfig {
	config := action.GruActionConfig{
		Service:         plan.Service,
		Target:          plan.Target,
		TargetType:      plan.TargetType,
		HostConfig:      &srv.HostConfig,
		ContainerConfig: &srv.ContainerConfig,
	}

	return &config
}
