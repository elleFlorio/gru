package strategy

import (
	"github.com/elleFlorio/gru/data"
)

type probDeltaStrategy struct{}

func (p *probDeltaStrategy) Name() string {
	return "probdelta"
}

func (p *probDeltaStrategy) Initialize() error {
	return nil
}

func (p *probDeltaStrategy) MakeDecision(policies []data.Policy) *data.Policy {
	threshold := randUniform(0, 1)
	shuffle(policies)
	return deltaElement(policies, threshold)
}

func deltaElement(policies []data.Policy, threshold float64) *data.Policy {
	var chosenPolicy data.Policy
	totalWeight := 0.0
	delta := 1.0

	for _, plc := range policies {
		totalWeight += plc.Weight
	}

	wNorm := 0.0
	wDelta := 0.0
	for _, plc := range policies {
		wNorm = plc.Weight / totalWeight
		if wNorm > threshold {
			chosenPolicy = plc
			return &chosenPolicy
		} else {
			wDelta = threshold - wNorm
			if wDelta < delta {
				delta = wDelta
				chosenPolicy = plc
			}
		}
	}

	return &chosenPolicy
}
