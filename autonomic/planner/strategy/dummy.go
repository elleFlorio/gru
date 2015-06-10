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

	log.WithFields(log.Fields{
		"status": "plan chosen",
		"plan":   thePlan,
	}).Debugln("Making decision")

	return thePlan, err
}

func (p *DummyStrategy) chosePlan(plans []GruPlan) *GruPlan {
	weight := 0.0
	thePlan := GruPlan{
		Service:    "none",
		Weight:     weight,
		TargetType: "none",
		Target:     "none",
		Actions:    []string{"noAction"},
	}

	for _, plan := range plans {
		if plan.Weight > weight {
			thePlan = plan
			weight = thePlan.Weight
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
	case "none":
		target = "none"
	default:
		return "", ErrorNoSuchTarget
	}

	return target, nil
}
