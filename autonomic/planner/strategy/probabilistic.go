package strategy

import (
	"errors"
	"math/rand"
	"time"

	"github.com/elleFlorio/gru/data"
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

func (p *probabilisticStrategy) MakeDecision(policies []data.Policy) *data.Policy {
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

func randUniform(min, max float64) float64 {
	return gen.Float64()*(max-min) + min
}

func shuffle(policies []data.Policy) {
	for i := range policies {
		j := gen.Intn(i + 1)
		policies[i], policies[j] = policies[j], policies[i]
	}
}
