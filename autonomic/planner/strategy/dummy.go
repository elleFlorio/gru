package strategy

import (
	"github.com/elleFlorio/gru/autonomic/planner/policy"
)

type dummyStrategy struct{}

func (p *dummyStrategy) Name() string {
	return "dummy"
}

func (p *dummyStrategy) Initialize() error {
	return nil
}

func (p *dummyStrategy) MakeDecision(policies []policy.Policy) *policy.Policy {
	var chosenPolicy *policy.Policy
	maxWeight := 0.0
	for _, plc := range policies {
		if plc.Weight > maxWeight {
			chosenPolicy = &plc
			maxWeight = plc.Weight
		}
	}

	return chosenPolicy
}
