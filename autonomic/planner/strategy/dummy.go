package strategy

import (
	"github.com/elleFlorio/gru/enum"
)

type dummyStrategy struct{}

func (p *dummyStrategy) Name() string {
	return "dummy"
}

func (p *dummyStrategy) Initialize() error {
	return nil
}

func (p *dummyStrategy) MakeDecision(plans []GruPlan) *GruPlan {
	thePlan := createNoActionPlan()
	maxWeight := enum.ValueFrom(enum.WHITE)
	for _, plan := range plans {
		if enum.ValueFrom(plan.Label) > maxWeight {
			thePlan = plan
			maxWeight = enum.ValueFrom(plan.Label)
		}
	}
	return &thePlan
}
