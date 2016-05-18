package strategy

import (
	"math/rand"
	"time"

	"github.com/elleFlorio/gru/data"
)

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	gen = rand.New(source)
}

var (
	gen *rand.Rand
)

func randUniform(min, max float64) float64 {
	return gen.Float64()*(max-min) + min
}

func shuffle(policies []data.Policy) {
	for i := range policies {
		j := gen.Intn(i + 1)
		policies[i], policies[j] = policies[j], policies[i]
	}
}
