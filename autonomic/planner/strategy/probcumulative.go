package strategy

import (
	"github.com/elleFlorio/gru/data"
)

type probCumulativeStrategy struct{}

func (p *probCumulativeStrategy) Name() string {
	return "probcumulative"
}

func (p *probCumulativeStrategy) Initialize() error {
	return nil
}

func (p *probCumulativeStrategy) MakeDecision(policies []data.Policy) *data.Policy {
	threshold := randUniform(0, 1)
	shuffle(policies)
	return weightedRandomElement(policies, threshold)
}

func weightedRandomElement(policies []data.Policy, threshold float64) *data.Policy {
	var chosenPolicy *data.Policy
	totalWeight := 0.0
	normalizedCumulative := 0.0

	for _, plc := range policies {
		totalWeight += plc.Weight
	}

	for _, plc := range policies {
		normalizedCumulative += plc.Weight / totalWeight
		if normalizedCumulative > threshold {
			chosenPolicy = &plc
			return chosenPolicy
		}
	}

	return chosenPolicy

}
