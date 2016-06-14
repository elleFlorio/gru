package metric

import (
	"math"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/data"
	"github.com/elleFlorio/gru/enum"
	res "github.com/elleFlorio/gru/resources"
	srv "github.com/elleFlorio/gru/service"
)

type updatePackage struct {
	target     string
	metricType enum.MetricType
	metric     string
	toAdd      []float64
}

var (
	servicesMetrics  map[string]Metric
	instancesMetrics map[string]Metric
	ch_inst_add      chan string
	ch_inst_rm       chan string
	ch_update        chan updatePackage
	ch_compute       chan struct{}
	ch_metrics       chan data.MetricStats
)

func init() {
	ch_inst_add = make(chan string, 100)
	ch_inst_rm = make(chan string, 100)
	ch_update = make(chan updatePackage, 100)
	ch_compute = make(chan struct{})
	ch_metrics = make(chan data.MetricStats)

	instancesMetrics = make(map[string]Metric)
}

func Initialize() {
	servicesList := srv.List()
	servicesMetrics = make(map[string]Metric, len(servicesList))
	for _, service := range servicesList {
		metric := Metric{
			UserMetrics: make(map[string][]float64),
		}

		servicesMetrics[service] = metric
	}
}

func AddInstance(id string) {
	ch_inst_add <- id
}

func RemoveInstance(id string) {
	ch_inst_rm <- id
}

func UpdateCpuMetric(id string, toAddInst []float64, toAddSys []float64) {
	ch_update <- updatePackage{
		target:     id,
		metricType: enum.METRIC_T_BASE,
		metric:     enum.METRIC_CPU_INST.ToString(),
		toAdd:      toAddInst,
	}

	ch_update <- updatePackage{
		target:     id,
		metricType: enum.METRIC_T_BASE,
		metric:     enum.METRIC_CPU_SYS.ToString(),
		toAdd:      toAddSys,
	}
}

func UpdateMemMetric(id string, toAdd []float64) {
	ch_update <- updatePackage{
		target:     id,
		metricType: enum.METRIC_T_BASE,
		metric:     enum.METRIC_MEM_INST.ToString(),
		toAdd:      toAdd,
	}
}

func UpdateUserMetric(service string, metric string, toAdd []float64) {
	ch_update <- updatePackage{
		target:     service,
		metricType: enum.METRIC_T_USER,
		metric:     metric,
		toAdd:      toAdd,
	}
}

func GetMetricsStats() data.MetricStats {
	ch_compute <- struct{}{}
	metStats := <-ch_metrics

	return metStats
}

func StartMetricCollector() {
	go metricCollector()
}

func metricCollector() {
	for {
		select {
		case id := <-ch_inst_add:
			instancesMetrics[id] = Metric{
				BaseMetrics: make(map[string][]float64),
			}
		case id := <-ch_inst_rm:
			delete(instancesMetrics, id)
		case update := <-ch_update:
			updateMetric(update.target, update.metricType, update.metric, update.toAdd)
		case <-ch_compute:
			metricsStats := computeMetrics()
			ch_metrics <- metricsStats
			clearMetrics()
		}
	}
}

func updateMetric(target string, metricType enum.MetricType, metric string, toAdd []float64) {

	switch metricType {
	case enum.METRIC_T_BASE:
		if toUpdateInstace, ok := instancesMetrics[target]; ok {
			if values, ok := toUpdateInstace.BaseMetrics[metric]; ok {
				values = append(values, toAdd...)
				toUpdateInstace.BaseMetrics[metric] = values
			} else {
				log.WithFields(log.Fields{
					"target": target,
					"metric": metric,
				}).Errorln("Cannot update instance metric: unknown metric")
			}
		} else {
			log.WithFields(log.Fields{
				"target": target,
				"metric": metric,
			}).Errorln("Cannot update instance metric: unknown instance")
		}
	case enum.METRIC_T_USER:
		if toUpdateService, ok := servicesMetrics[target]; ok {
			if values, ok := toUpdateService.UserMetrics[metric]; ok {
				values = append(values, toAdd...)
				toUpdateService.UserMetrics[metric] = values
			} else {
				log.WithFields(log.Fields{
					"target": target,
					"metric": metric,
				}).Errorln("Cannot update service metric: unknown metric")
			}
		} else {
			log.WithFields(log.Fields{
				"target": target,
				"metric": metric,
			}).Errorln("Cannot update service metric: unknown instance")
		}
	}
}

func computeMetrics() data.MetricStats {
	metrics := data.MetricStats{}
	instMetrics := computeInstancesMetrics()
	servMetrics := computeServicesMetrics(instMetrics)
	sysMetrics := computeSysMetrics(servMetrics)

	metrics.Instance = instMetrics
	metrics.Service = servMetrics
	metrics.System = sysMetrics

	return metrics
}

func computeInstancesMetrics() map[string]data.MetricData {
	instMetrics := make(map[string]data.MetricData)
	baseMetrics := make(map[string]float64)

	for instance, metrics := range instancesMetrics {
		// CPU
		instCpus := metrics.BaseMetrics[enum.METRIC_CPU_INST.ToString()]
		sysCpus := metrics.BaseMetrics[enum.METRIC_CPU_SYS.ToString()]
		cpuPerc := computeInstanceCpuPerc(instCpus, sysCpus)
		baseMetrics[enum.METRIC_CPU_AVG.ToString()] = cpuPerc

		// MEMORY - TODO
		memPerc := 0.0
		baseMetrics[enum.METRIC_MEM_AVG.ToString()] = memPerc

		instMetrics[instance] = data.MetricData{
			BaseMetrics: baseMetrics,
		}
	}

	return instMetrics
}

// Since linux compute the cpu usage in units of jiffies, it needs to be converted
// in % using the formula used in this function.
// Explaination: http://stackoverflow.com/questions/1420426/calculating-cpu-usage-of-a-process-in-linux
// TODO probably I just need the first and the last value...
// 2015/11/16 - corrected according to what the docker client does:
// https://github.com/docker/docker/blob/master/api/client/stats.go#L316
func computeInstanceCpuPerc(instCpus []float64, sysCpus []float64) float64 {
	sum := 0.0
	instNext := 0.0
	sysNext := 0.0
	instPrev := 0.0
	sysPrev := 0.0
	cpu := 0.0
	cpuTotal := res.GetResources().CPU.Total

	valid := 0
	nValues := int(math.Min(float64(len(instCpus)), float64(len(sysCpus))))

	for i := 1; i < nValues; i++ {
		instPrev = instCpus[i-1]
		sysPrev = sysCpus[i-1]
		instNext = instCpus[i]
		sysNext = sysCpus[i]
		instDelta := instNext - instPrev
		if instDelta > 0 {
			sysDelta := sysNext - sysPrev
			if sysDelta == 0 {
				cpu = 0
			} else {
				// "100 * cpu" should produce values in [0, 100]
				cpu = (instDelta / sysDelta) * float64(cpuTotal)
			}
			sum += cpu
			valid++
		}
	}

	if valid > 0.0 {
		return math.Min(1.0, sum/float64(valid))
	}

	return 0.0
}

func computeServicesMetrics(instMetrics map[string]data.MetricData) map[string]data.MetricData {

	servicesAvg := make(map[string]data.MetricData, len(servicesMetrics))

	for service, metrics := range servicesMetrics {
		baseMetrics := make(map[string]float64)
		// CPU
		cpuAvg, cpuTot := computeServiceCpuPerc(service, instMetrics)
		baseMetrics[enum.METRIC_CPU_AVG.ToString()] = cpuAvg
		baseMetrics[enum.METRIC_CPU_TOT.ToString()] = cpuTot
		// MEMORY -TODO
		memAvg := 0.0
		memTot := 0.0
		baseMetrics[enum.METRIC_MEM_AVG.ToString()] = memAvg
		baseMetrics[enum.METRIC_MEM_TOT.ToString()] = memTot

		userMetrics := make(map[string]float64, len(metrics.UserMetrics))
		for metric, values := range metrics.UserMetrics {
			value := mean(values)
			userMetrics[metric] = value
		}

		serviceAvg := data.MetricData{
			BaseMetrics: baseMetrics,
			UserMetrics: userMetrics,
		}

		servicesAvg[service] = serviceAvg
	}

	return servicesAvg
}

// Returns CPU percentage average, total.
func computeServiceCpuPerc(name string, instMetrics map[string]data.MetricData) (float64, float64) {
	sum := 0.0

	service, _ := srv.GetServiceByName(name)
	values := make([]float64, 0)

	if len(service.Instances.Running) > 0 {
		for _, id := range service.Instances.Running {
			instCpuAvg := instMetrics[id].BaseMetrics[enum.METRIC_CPU_AVG.ToString()]
			sum += instCpuAvg
			values = append(values, instCpuAvg)
		}
	}

	return mean(values), math.Min(1.0, sum)
}

func mean(values []float64) float64 {
	if len(values) < 1 {
		return 0.0
	}

	sum := 0.0
	for _, value := range values {
		sum += value
	}

	return sum / float64(len(values))
}

func computeSysMetrics(servMetrics map[string]data.MetricData) data.MetricData {
	// TODO - improve by adding capacity
	baseMetrics := make(map[string]float64)
	cpuSys := 0.0
	memSys := 0.0
	for _, metrics := range servMetrics {
		// CPU
		cpuSys += metrics.BaseMetrics[enum.METRIC_CPU_TOT.ToString()]

		// MEM
		// TODO
	}

	baseMetrics[enum.METRIC_CPU_TOT.ToString()] = cpuSys
	baseMetrics[enum.METRIC_MEM_TOT.ToString()] = memSys
	sysMetrics := data.MetricData{
		BaseMetrics: baseMetrics,
	}

	return sysMetrics

}

// FIXME
func clearMetrics() {
	capacity := len(servicesMetrics)
	servicesMetrics = make(map[string]Metric, capacity)
}
