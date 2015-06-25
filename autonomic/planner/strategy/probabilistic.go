package strategy

import (
	"errors"
	"math/rand"

	log "github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/autonomic/analyzer"
	"github.com/elleFlorio/gru/service"
)

var (
	ErrorTotalWeightIsZero error = errors.New("Total weight of plans is zero")
)

type ProbabilisticStrategy struct{}

func (p *ProbabilisticStrategy) Name() string {
	return "probabilistic"
}

func (p *ProbabilisticStrategy) Initialize() error {
	return nil
}

func (p *ProbabilisticStrategy) MakeDecision(plans []GruPlan, analytics *analyzer.GruAnalytics) (*GruPlan, error) {
	thePlan := p.chosePlan(plans)
	srv, _ := service.GetServiceByName(thePlan.Service)
	target, err := p.choseTarget(thePlan.TargetType, thePlan.TargetStatus, analytics, srv)

	// FIXME this is not good. Find a way to do it better
	// If I don't have stopped containers to start,
	// I have to create a new one starting from an image
	if err == ErrorNoStoppedCont {
		thePlan.TargetType = "image"
		target = srv.Image
	}

	thePlan.Target = target

	log.WithFields(log.Fields{
		"status": "plan chosen",
		"plan":   thePlan,
	}).Debugln("Making decision")

	return thePlan, err
}

func (p *ProbabilisticStrategy) chosePlan(plans []GruPlan) *GruPlan {
	thePlan, err := p.weightedRandomElement(plans)
	if err != nil {
		log.WithFields(log.Fields{
			"status": "chosing plan",
			"err":    err,
		}).Debugln("Making decision")

		thePlan = &GruPlan{
			Service:    "none",
			Weight:     0.0,
			TargetType: "none",
			Target:     "none",
			Actions:    []string{"noAction"},
		}
	}

	return thePlan

}

func (p *ProbabilisticStrategy) weightedRandomElement(plans []GruPlan) (*GruPlan, error) {
	totalWeight := 0.0
	threshold := p.randUniform(0, 1)
	normalizedCumulative := 0.0

	for _, plan := range plans {
		totalWeight += plan.Weight
	}

	for _, plan := range plans {
		normalizedCumulative += plan.Weight / totalWeight
		if normalizedCumulative > threshold {
			return &plan, nil
		}
	}

	return nil, ErrorTotalWeightIsZero

}

func (p *ProbabilisticStrategy) randUniform(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

func (p *ProbabilisticStrategy) choseTarget(tType string, tStatus string, analytics *analyzer.GruAnalytics, srv *service.Service) (string, error) {
	var target string
	var pool []string
	switch tType {
	case "container":
		instances := analytics.Service[srv.Name].Instances
		switch tStatus {
		case "running":
			pool = instances.Active
		case "stopped":
			if len(instances.Stopped) > 0 {
				pool = instances.Stopped
			} else {
				return "", ErrorNoStoppedCont
			}
		default:
			return "", ErrorNoSuchStatus
		}
		target = pool[rand.Intn(len(pool))]
	case "image":
		target = srv.Image
	case "none":
		target = "none"
	default:
		return "", ErrorNoSuchTarget
	}

	return target, nil
}
