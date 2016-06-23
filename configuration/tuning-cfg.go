package configuration

type Tuning struct {
	Policy PolicyThreshold
}

type PolicyThreshold struct {
	Scaleout float64
	Scalein  float64
	Swap     float64
}
