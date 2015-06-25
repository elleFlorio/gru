package strategy

func CreateMockPlans(w1, w2, w3 float64) []GruPlan {
	p1 := GruPlan{
		Service:    "service1",
		Weight:     w1,
		TargetType: "container",
		Actions:    []string{"start", "stop"},
	}

	p2 := GruPlan{
		Service:    "service2",
		Weight:     w2,
		TargetType: "image",
		Actions:    []string{"open"},
	}

	p3 := GruPlan{
		Service:    "service3",
		Weight:     w3,
		TargetType: "notExist",
		Actions:    []string{"close, shutdown"},
	}

	return []GruPlan{p1, p2, p3}

}
