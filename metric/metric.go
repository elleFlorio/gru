package metric

import (
	"github.com/elleFlorio/gru/data"
)

type GruMetric struct {
	Node    NodeMetrics
	Service map[string]ServiceMetric
	Policy  PolicyMetric
}

type NodeMetrics struct {
	UUID           string
	Name           string
	Stats          data.MetricData
	Resources      ResourceMetrics
	ActiveServices int
}

type ResourceMetrics struct {
	CPU CpuMetrics
}

type CpuMetrics struct {
	Availabe int64
	Total    int64
}

type ServiceMetric struct {
	Name      string
	Type      string
	Image     string
	Instances InstancesMetric
	Stats     data.MetricData
	Analytics data.AnalyticData
	Shared    data.SharedData
}

type InstancesMetric struct {
	All     int
	Pending int
	Running int
	Stopped int
	Paused  int
}

type PolicyMetric struct {
	Name   string
	Weight float64
}
