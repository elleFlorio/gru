package executor

import (
	log "github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/action"
	"github.com/elleFlorio/gru/autonomic/planner/strategy"
)

type executor struct{}

func NewExecutor() *executor {
	return &executor{}
}

func (p *executor) Run(plan strategy.GruPlan) {
	//Execute stuff
	log.Debugln("I'm executing")
	actions := buildActions(&plan)
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
