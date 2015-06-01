package strategy

import (
	"math/rand"

	"github.com/elleFlorio/gru/autonomic/analyzer"
	"github.com/elleFlorio/gru/service"
)

type DummyStrategy struct{}

func (p *DummyStrategy) Name() string {
	return "dummy"
}

func (p *DummyStrategy) Initialize() error {
	return nil
}

func (p *DummyStrategy) MakeDecision(plans []GruPlan, analytics *analyzer.GruAnalytics) (*GruPlan, error) {
	thePlan := p.chosePlan(plans)
	err := p.choseTarget(thePlan, analytics)
	return thePlan, err
}

func (p *DummyStrategy) chosePlan(plans []GruPlan) *GruPlan {
	var thePlan *GruPlan
	weigth := 0.0
	for _, plan := range plans {
		if plan.Weight > weigth {
			thePlan = &plan
		}
	}

	return thePlan
}

func (p *DummyStrategy) choseTarget(thePlan *GruPlan, analytics *analyzer.GruAnalytics) error {
	var target string
	switch thePlan.TargetType {
	case "container":
		instances := analytics.Service[thePlan.Service].Instances
		target = instances[rand.Intn(len(instances))]
	case "image":
		srv, _ := service.GetServiceByName(thePlan.Service)
		target = srv.Image
	default:
		return ErrorNoSuchTarget
	}
	thePlan.Target = target
	return nil
}
