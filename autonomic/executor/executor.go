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
	actions, err := p.buildActions(&plan)
	if err != nil {
		p.c_err <- err
	}
	srv, _ := service.GetServiceByName(plan.Service)
	config := p.buildConfig(&plan, srv, docker)
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

func (p *executor) buildConfig(plan *strategy.GruPlan, srv *service.Service, docker *dockerclient.DockerClient) *action.GruActionConfig {
	config := action.GruActionConfig{
		Service:    plan.Service,
		Target:     plan.Target,
		Client:     docker,
		HostConfig: &srv.HostConfig,
	}

	return &config
}
