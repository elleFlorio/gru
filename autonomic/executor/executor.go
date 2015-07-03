package executor

import (
	log "github.com/Sirupsen/logrus"
	"github.com/samalba/dockerclient"

	"github.com/elleFlorio/gru/action"
	"github.com/elleFlorio/gru/autonomic/planner/strategy"
	"github.com/elleFlorio/gru/service"
)

type executor struct{ c_err chan error }

func NewExecutor(c_err chan error) *executor {
	return &executor{c_err}
}

func (p *executor) Run(plan strategy.GruPlan, docker *dockerclient.DockerClient) {
	log.WithField("status", "start").Debugln("Running executor")
	defer log.WithField("status", "done").Debugln("Running executor")

	config := &action.GruActionConfig{}
	actions, err := p.buildActions(&plan)
	if err != nil {
		log.WithFields(log.Fields{
			"status": "error",
			"error":  err,
		}).Debugln("Running executor")
	}

	srv, err := service.GetServiceByName(plan.Service)

	if err == nil {
		config = p.buildConfig(&plan, srv, docker)
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

func (p *executor) buildActions(plan *strategy.GruPlan) ([]action.GruAction, error) {
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
func (p *executor) buildConfig(plan *strategy.GruPlan, srv *service.Service, docker *dockerclient.DockerClient) *action.GruActionConfig {
	config := action.GruActionConfig{
		Service:         plan.Service,
		Target:          plan.Target,
		TargetType:      plan.TargetType,
		Client:          docker,
		HostConfig:      &srv.HostConfig,
		ContainerConfig: &srv.ContainerConfig,
	}

	return &config
}
