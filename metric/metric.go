package metric

type GruMetric struct {
	Node    NodeMetrics
	Service map[string]ServiceMetric
	Plan    PlansMetric
}

type NodeMetrics struct {
	UUID   string
	Name   string
	Cpu    float64
	Memory float64
	Health float64
}

type ServiceMetric struct {
	Name      string
	Type      string
	Image     string
	Instances InstancesMetric
	Stats     StatsMetric
	Analytics AnalyticsMetric
}

type InstancesMetric struct {
	All     int
	Pending int
	Running int
	Stopped int
	Paused  int
}

type StatsMetric struct {
	CpuAvg float64
	CpuTot float64
	MemAvg float64
	MemTot float64
}

type AnalyticsMetric struct {
	Cpu       float64
	Memory    float64
	Resources float64
	Load      float64
	Health    float64
}

type PlansMetric struct {
	Policy string
	Target string
	Weight float64
}
