package metric

type ServiceMetric struct {
	UserMetric map[string][]float64
}

type InstanceMetric struct {
	BaseMetric map[string][]float64
}
