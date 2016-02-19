package policy

import (
	"math/rand"
	"time"

	"github.com/elleFlorio/gru/enum"
)

const c_LETTERBYTES = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var gen *rand.Rand

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	gen = rand.New(source)
}

func CreateMockPolicy(name string, weight float64, targets map[string][]enum.Action) Policy {
	return Policy{name, weight, targets}
}

func CreateRandomMockPolicies(nServices int) []Policy {
	srvList := make([]string, nServices, nServices)
	for i := 0; i < nServices; i++ {
		name := randStringBytes(5)
		srvList[i] = name
	}

	return createRandomScalePolicies(srvList)
}

func createRandomScalePolicies(srvList []string) []Policy {
	policies := make([]Policy, 0, len(srvList))
	scale := []string{"scaleout", "scalein"}
	for _, inOut := range scale {
		for _, srv := range srvList {
			name := inOut
			weight := randUniform(0, 1)
			targets := map[string][]enum.Action{srv: []enum.Action{enum.STOP}}
			policies = append(policies, CreateMockPolicy(name, weight, targets))
		}
	}

	return policies
}

func randUniform(min, max float64) float64 {
	return gen.Float64()*(max-min) + min
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = c_LETTERBYTES[rand.Intn(len(c_LETTERBYTES))]
	}
	return string(b)
}
