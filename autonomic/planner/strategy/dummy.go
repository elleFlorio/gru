package strategy

type dummyStrategy struct{}

func (p *dummyStrategy) Name() string {
	return "dummy"
}

func (p *dummyStrategy) Initialize() error {
	return nil
}

func (p *dummyStrategy) MakeDecision(plans []GruPlan) *GruPlan {
	var thePlan GruPlan
	maxWeight := 0.0
	for _, plan := range plans {
		if plan.Weight > maxWeight {
			thePlan = plan
			maxWeight = plan.Weight
		}
	}
	return &thePlan
}
