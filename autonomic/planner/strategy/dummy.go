package strategy

import (
	"math/rand"

	log "github.com/Sirupsen/logrus"

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
	srv, _ := service.GetServiceByName(thePlan.Service)
	target, err := p.choseTarget(thePlan.TargetType, analytics, srv)
	if err != nil {
		log.WithFields(log.Fields{
			"status": "target",
			"error":  err,
		}).Errorln("Making decision")
	}
	thePlan.Target = target
	return thePlan, err
}

func (p *DummyStrategy) chosePlan(plans []GruPlan) *GruPlan {
	var thePlan GruPlan
	weigth := 0.0
	for _, plan := range plans {
		if plan.Weight > weigth {
			thePlan = plan
			weigth = thePlan.Weight
		}
	}

	return &thePlan
}

func (p *DummyStrategy) choseTarget(tType string, analytics *analyzer.GruAnalytics, srv *service.Service) (string, error) {
	var target string
	switch tType {
	case "container":
		instances := analytics.Service[srv.Name].Instances
		target = instances[rand.Intn(len(instances))]
	case "image":
		target = srv.Image
	default:
		return "", ErrorNoSuchTarget
	}

	return target, nil
}
