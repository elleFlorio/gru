package metric

import (
	"math"
	"sync"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/data"
	"github.com/elleFlorio/gru/enum"
	res "github.com/elleFlorio/gru/resources"
	srv "github.com/elleFlorio/gru/service"
	"github.com/elleFlorio/gru/utils"
)

var (
	servicesMetrics  map[string]Metric
	instancesMetrics map[string]Metric

	mutex_instMet sync.RWMutex
)

func init() {
	instancesMetrics = make(map[string]Metric)
	mutex_instMet = sync.RWMutex{}
}

func Initialize(services []string) {
	servicesMetrics = make(map[string]Metric, len(services))
	for _, service := range services {
		metric := Metric{
			UserMetrics: make(map[string][]float64),
		}

		servicesMetrics[service] = metric
	}
}

func AddInstance(id string) {
	defer mutex_instMet.Unlock()

	mutex_instMet.Lock()
	instancesMetrics[id] = Metric{
		BaseMetrics: make(map[string][]float64),
	}

	log.WithField("id", id).Debugln("Added instance to metric collector")
}

func RemoveInstance(id string) {
	defer mutex_instMet.Unlock()

	mutex_instMet.Lock()
	delete(instancesMetrics, id)

	log.WithField("id", id).Debugln("Removed instance from metric collector")
}

func UpdateCpuMetric(id string, toAddInst []float64, toAddSys []float64) {
	updateMetric(id, enum.METRIC_T_BASE, enum.METRIC_CPU_INST.ToString(), toAddInst)
	updateMetric(id, enum.METRIC_T_BASE, enum.METRIC_CPU_SYS.ToString(), toAddSys)
}

func UpdateMemMetric(id string, toAdd []float64) {
	updateMetric(id, enum.METRIC_T_BASE, enum.METRIC_MEM_INST.ToString(), toAdd)
}

func UpdateUserMetric(service string, metric string, toAdd []float64) {
	updateMetric(service, enum.METRIC_T_USER, metric, toAdd)
}

func IsReadyForRunning(instance string, threshold int) bool {
	defer mutex_instMet.RUnlock()

	mutex_instMet.RLock()
	readyToRun := true
	metrics := instancesMetrics[instance].BaseMetrics
	for _, values := range metrics {
		readyToRun = readyToRun && (len(values) >= threshold)
	}

	return readyToRun
}

func GetMetricsStats() data.MetricStats {
	defer resetMetrics()
	metStats := computeMetrics()

	return metStats
}

func updateMetric(target string, metricType enum.MetricType, metric string, toAdd []float64) {
	defer mutex_instMet.Unlock()

	mutex_instMet.Lock()
	switch metricType {
	case enum.METRIC_T_BASE:
		if toUpdateInstace, ok := instancesMetrics[target]; ok {
			values := toUpdateInstace.BaseMetrics[metric]
			values = append(values, toAdd...)
			toUpdateInstace.BaseMetrics[metric] = values
		} else {
			log.WithFields(log.Fields{
				"target": target,
				"metric": metric,
			}).Errorln("Cannot update instance metric: unknown instance")
		}
	case enum.METRIC_T_USER:
		if toUpdateService, ok := servicesMetrics[target]; ok {
			values := toUpdateService.UserMetrics[metric]
			values = append(values, toAdd...)
			toUpdateService.UserMetrics[metric] = values
		} else {
			log.WithFields(log.Fields{
				"target": target,
				"metric": metric,
			}).Errorln("Cannot update service metric: unknown service")
		}
	}
}

func computeMetrics() data.MetricStats {
	metrics := data.MetricStats{}
	instMetrics := computeInstancesMetrics()
	log.WithField("instMetrics", instMetrics).Debugln("Computed instances metrics")
	servMetrics := computeServicesMetrics(instMetrics)
	log.WithField("servMetrics", servMetrics).Debugln("Computed service metrics")
	sysMetrics := computeSysMetrics(instMetrics)
	log.WithField("sysMetrics", sysMetrics).Debugln("Computed system metrics")

	metrics.Instance = instMetrics
	metrics.Service = servMetrics
	metrics.System = sysMetrics

	return metrics
}

func computeInstancesMetrics() map[string]data.MetricData {
	defer mutex_instMet.RUnlock()

	mutex_instMet.RLock()
	instMetrics := make(map[string]data.MetricData)

	for instance, metrics := range instancesMetrics {
		baseMetrics := make(map[string]float64)

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
		cpuAvg := computeServiceCpuPerc(service, instMetrics)
		baseMetrics[enum.METRIC_CPU_AVG.ToString()] = cpuAvg
		// MEMORY -TODO
		memAvg := 0.0
		baseMetrics[enum.METRIC_MEM_AVG.ToString()] = memAvg

		userMetrics := make(map[string]float64, len(metrics.UserMetrics))
		for metric, values := range metrics.UserMetrics {
			value := utils.Mean(values)
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
func computeServiceCpuPerc(name string, instMetrics map[string]data.MetricData) float64 {

	service, _ := srv.GetServiceByName(name)
	values := make([]float64, 0)

	if len(service.Instances.Running) > 0 {
		for _, id := range service.Instances.Running {
			instCpuAvg := instMetrics[id].BaseMetrics[enum.METRIC_CPU_AVG.ToString()]
			values = append(values, instCpuAvg)
		}
	}

	return utils.Mean(values)
}

func computeSysMetrics(instMetrics map[string]data.MetricData) data.MetricData {
	// TODO - improve by adding capacity
	baseMetrics := make(map[string]float64)
	cpuSys := 0.0
	memSys := make([]float64, 0, len(instMetrics))
	for instance, metrics := range instMetrics {
		service, err := srv.GetServiceById(instance)
		if err != nil {
			log.WithFields(log.Fields{
				"instance": instance,
			}).Errorln("Cannot find service by instance")
		} else {
			instCpus := service.Docker.CPUnumber
			instCpuValue := metrics.BaseMetrics[enum.METRIC_CPU_AVG.ToString()] * float64(instCpus)
			// CPU
			cpuSys += instCpuValue

			// MEM
			// TODO
		}
	}

	baseMetrics[enum.METRIC_CPU_AVG.ToString()] = cpuSys / float64(res.GetResources().CPU.Total)
	baseMetrics[enum.METRIC_MEM_AVG.ToString()] = utils.Mean(memSys)
	sysMetrics := data.MetricData{
		BaseMetrics: baseMetrics,
	}

	return sysMetrics

}

func resetMetrics() {
	defer mutex_instMet.Unlock()

	mutex_instMet.Lock()
	for service, metrics := range servicesMetrics {
		for metric, _ := range metrics.UserMetrics {
			metrics.UserMetrics[metric] = metrics.UserMetrics[metric][:0]
		}
		servicesMetrics[service] = metrics
	}

	for instance, metrics := range instancesMetrics {
		for metric, _ := range metrics.BaseMetrics {
			metrics.BaseMetrics[metric] = metrics.BaseMetrics[metric][:0]
		}
		instancesMetrics[instance] = metrics
	}
}
