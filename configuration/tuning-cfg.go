package configuration

type Tuning struct {
	Policy PolicyThreshold
}

type PolicyThreshold struct {
	Scaleout ScaleThreshold
	Scalein  ScaleThreshold
	Swap     ScaleThreshold
}

type ScaleThreshold struct {
	Cpu  float64
	Load float64
}
