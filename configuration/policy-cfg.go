package configuration

type Policy struct {
	Scalein  PolicyConfig
	Scaleout PolicyConfig
	Swap     PolicyConfig
}

type PolicyConfig struct {
	Enable    bool
	Threshold float64
	Metrics   []string
	Analytics []string
}
