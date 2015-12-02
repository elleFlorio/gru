package metric

type metricService interface {
	Name() string
	Initialize(interface{}) string
	StoreMetrics(GruMetric) error
	GetMetrics(interface{}) (GruMetric, error)
}
