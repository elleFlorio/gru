package strategy

import (
	"math/rand"

	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/service"
	"github.com/elleFlorio/gru/storage"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func CreateMockPlan(l enum.Label, s service.Service, a []enum.Action) GruPlan {
	return GruPlan{l, &s, a}
}

func StoreMockPlan(l enum.Label, s service.Service, a []enum.Action) {
	plan := CreateMockPlan(l, s, a)
	data, _ := ConvertPlanToData(plan)
	storage.StoreLocalData(data, enum.PLANS)
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
