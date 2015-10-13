package strategy

import (
	"errors"
	"math/rand"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/enum"
)

var (
	ErrorTotalWeightIsZero   error = errors.New("Total weight of plans is zero")
	ErrorThresholdNotReached error = errors.New("Threshold not reached")
)

type probabilisticStrategy struct{}

func (p *probabilisticStrategy) Name() string {
	return "probabilistic"
}

func (p *probabilisticStrategy) Initialize() error {
	return nil
}

func (p *probabilisticStrategy) MakeDecision(plans []GruPlan) *GruPlan {
	thePlan, err := weightedRandomElement(plans)
	if err != nil {
		log.WithField("err", err).Debugln("returning no action")
		return &GruPlan{Actions: []enum.Action{enum.NOACTION}}
	}

	return thePlan
}

func weightedRandomElement(plans []GruPlan) (*GruPlan, error) {
	totalWeight := 0.0
	threshold := randUniform(0, 1)
	normalizedCumulative := 0.0

	for _, plan := range plans {
		totalWeight += enum.ValueFrom(plan.Label)
	}

	if totalWeight == 0.0 {
		return nil, ErrorTotalWeightIsZero
	}

	shuffle(plans)

	for _, plan := range plans {
		normalizedCumulative += enum.ValueFrom(plan.Label) / totalWeight
		if normalizedCumulative > threshold {
			return &plan, nil
		}
	}

	return nil, ErrorThresholdNotReached

}

func randUniform(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

func shuffle(plans []GruPlan) {
	for i := range plans {
		j := rand.Intn(i + 1)
		plans[i], plans[j] = plans[j], plans[i]
	}
}
