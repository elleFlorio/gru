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
	return weightedRandomElement(policies)
}

func weightedRandomElement(policies []data.Policy) *data.Policy {
	var chosenPolicy *data.Policy
	totalWeight := 0.0
	threshold := randUniform(0, 1)
	normalizedCumulative := 0.0

	for _, plc := range policies {
		totalWeight += plc.Weight
	}

	shuffle(policies)

	for _, plc := range policies {
		normalizedCumulative += plc.Weight / totalWeight
		if normalizedCumulative > threshold {
			chosenPolicy = &plc
			return chosenPolicy
		}
	}

	return chosenPolicy

}
