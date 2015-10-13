package strategy

import (
	"math/rand"

	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/service"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func CreateMockPlan(l enum.Label, s service.Service, a []enum.Action) GruPlan {
	return GruPlan{l, &s, a}
}

func CreateRandomPlans(n int) []GruPlan {
	plans := []GruPlan{}
	for i := 0; i < n; i++ {
		value := randUniform(0, 1)
		l := enum.FromValue(value)
		s := service.Service{Name: randStringBytes(5)}
		a := []enum.Action{enum.START}
		if value > 0.5 {
			a = []enum.Action{enum.STOP}
		}
		plans = append(plans, GruPlan{l, &s, a})
	}

	return plans
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// func CreateMockPlans(w1, w2, w3 float64) []GruPlan {
// 	p1 := GruPlan{
// 		Service:    "service1",
// 		Weight:     w1,
// 		TargetType: "container",
// 		Actions:    []string{"start", "stop"},
// 	}

// 	p2 := GruPlan{
// 		Service:    "service2",
// 		Weight:     w2,
// 		TargetType: "image",
// 		Actions:    []string{"open"},
// 	}

// 	p3 := GruPlan{
// 		Service:    "service3",
// 		Weight:     w3,
// 		TargetType: "notExist",
// 		Actions:    []string{"close, shutdown"},
// 	}

// 	return []GruPlan{p1, p2, p3}

// }
