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
	//Execute stuff
	log.Debugln("I'm executing")
	actions := p.buildActions(&plan)
	config := p.buildConfig(&plan, docker)
	for _, act := range actions {
		act.Run(config)
	}
}

func (p *executor) buildActions(plan *strategy.GruPlan) []action.GruAction {
	actions := make([]action.GruAction, len(plan.Actions))
	for _, name := range plan.Actions {
		act, err := action.New(name)
		if err != nil {
			//TODO
		}
		actions = append(actions, act)
	}

	return actions
}

func (p *executor) buildConfig(plan *strategy.GruPlan, docker *dockerclient.DockerClient) *action.GruActionConfig {
	srv, _ := service.GetServiceByName(plan.Service)
	config := action.GruActionConfig{
		Service:    plan.Service,
		Target:     plan.Target,
		Client:     docker,
		HostConfig: &srv.HostConfig,
	}

	return &config
}
