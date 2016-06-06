package enum

type MetricType string
type Metric string

const (
	// type
	METRIC_T_BASE MetricType = "base_metric"
	METRIC_T_USER MetricType = "user_metric"
	// metric
	METRIC_CPU Metric = "cpu"
	METRIC_MEM Metric = "memory"
)

func (t MetricType) Value() float64 {
	var v float64
	switch {
	case t == METRIC_T_BASE:
		v = 0.0
	case t == METRIC_T_USER:
		v = 1.0
	}

	return v
}

func (m Metric) Value() float64 {
	var v float64
	switch {
	case m == METRIC_CPU:
		v = 0.0
	case m == METRIC_MEM:
		v = 1.0
	}

	return v
}

func (t MetricType) ToString() string {
	var s string
	switch {
	case t == METRIC_T_BASE:
		s = "base_metric"
	case t == METRIC_T_USER:
		s = "user_metric"
	}

	return s
}

func (m Metric) ToString() string {
	var s string
	switch {
	case m == METRIC_CPU:
		s = "cpu"
	case m == METRIC_MEM:
		s = "memory"
	}

	return s
}
