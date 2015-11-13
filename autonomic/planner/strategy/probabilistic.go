package strategy

import (
	"errors"
	"math/rand"
	"time"

	"github.com/elleFlorio/gru/enum"
)

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	gen = rand.New(source)
}

var (
	gen                      *rand.Rand
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
	return weightedRandomElement(plans)
}

func weightedRandomElement(plans []GruPlan) *GruPlan {
	var thePlan GruPlan
	totalWeight := 0.0
	threshold := randUniform(0, 1)
	normalizedCumulative := 0.0

	for _, plan := range plans {
		totalWeight += enum.ValueFrom(plan.Label)
	}

	shuffle(plans)

	for _, plan := range plans {
		normalizedCumulative += enum.ValueFrom(plan.Label) / totalWeight
		if normalizedCumulative > threshold {
			thePlan = plan
			break
		}
	}

	return &thePlan

}

func randUniform(min, max float64) float64 {
	return gen.Float64()*(max-min) + min
}

func shuffle(plans []GruPlan) {
	for i := range plans {
		j := gen.Intn(i + 1)
		plans[i], plans[j] = plans[j], plans[i]
	}
}
