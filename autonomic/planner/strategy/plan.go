package strategy

type GruPlan struct {
	Service      string
	Weight       float64
	TargetType   string
	TargetStatus string
	Target       string
	Actions      []string
}
