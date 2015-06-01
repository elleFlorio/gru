package planner

import (
	log "github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/autonomic/analyzer"
	"github.com/elleFlorio/gru/autonomic/planner/strategy"
	"github.com/elleFlorio/gru/policy"
	"github.com/elleFlorio/gru/service"
)

type planner struct {
	strategy strategy.GruStrategy
	c_err    chan error
}

func NewPlanner(strategyName string, c_err chan error) *planner {
	strategy, err := strategy.New(strategyName)
	if err != nil {
		//TODO
	}
	return &planner{
		strategy,
		c_err,
	}
}

func (p *planner) Run(analytics analyzer.GruAnalytics) strategy.GruPlan {
	//Plan stuff
	log.Debugln("I'm planning")
	plans := buildPlans(&analytics)
	thePlan, err := p.strategy.MakeDecision(plans, &analytics)
	if err != nil {
		//TODO
	}

	return *thePlan
}

func buildPlans(analytics *analyzer.GruAnalytics) []strategy.GruPlan {
	policies := policy.GetPolicies("proactive")
	plans := []strategy.GruPlan{}

	for _, name := range service.List() {
		for _, plc := range policies {
			srvc, err := service.GetServiceByName(name)
			if err != nil {
				//TODO
			}

			plan := strategy.GruPlan{
				Service:    name,
				Weight:     plc.Weight(srvc, analytics),
				TargetType: plc.Target(),
				Actions:    plc.Actions(),
			}

			plans = append(plans, plan)
		}
	}

	return plans
}
