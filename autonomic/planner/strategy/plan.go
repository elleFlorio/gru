package strategy

type GruPlan struct {
	Service    string
	Weight     float64
	TargetType string
	Target     string
	Actions    []string
}
