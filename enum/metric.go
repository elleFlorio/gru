package enum

type MetricType string
type Metric string

const (
	// type
	METRIC_T_BASE MetricType = "base_metric"
	METRIC_T_USER MetricType = "user_metric"
	// metric
	METRIC_CPU_INST Metric = "cpu_inst"
	METRIC_CPU_SYS  Metric = "cpu_sys"
	METRIC_CPU_AVG  Metric = "cpu_avg"
	METRIC_CPU_TOT  Metric = "cpu_tot"
	METRIC_MEM_INST Metric = "memory_inst"
	METRIC_MEM_SYS  Metric = "memory_sys"
	METRIC_MEM_AVG  Metric = "memory_avg"
	METRIC_MEM_TOT  Metric = "memory_tot"
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
	case m == METRIC_CPU_INST:
		v = 0.0
	case m == METRIC_CPU_SYS:
		v = 1.0
	case m == METRIC_CPU_AVG:
		v = 2.0
	case m == METRIC_CPU_TOT:
		v = 3.0
	case m == METRIC_MEM_INST:
		v = 4.0
	case m == METRIC_MEM_SYS:
		v = 5.0
	case m == METRIC_MEM_AVG:
		v = 6.0
	case m == METRIC_MEM_TOT:
		v = 7.0
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
	case m == METRIC_CPU_INST:
		s = "cpu_inst"
	case m == METRIC_CPU_SYS:
		s = "cpu_sys"
	case m == METRIC_CPU_AVG:
		s = "cpu_avg"
	case m == METRIC_CPU_TOT:
		s = "cpu_tot"
	case m == METRIC_MEM_INST:
		s = "memory_inst"
	case m == METRIC_MEM_SYS:
		s = "memory_sys"
	case m == METRIC_MEM_AVG:
		s = "memory_avg"
	case m == METRIC_MEM_TOT:
		s = "memory_tot"
	}

	return s
}
