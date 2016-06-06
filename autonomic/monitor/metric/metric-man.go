package metric

import (
	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/data"
	"github.com/elleFlorio/gru/enum"
	srv "github.com/elleFlorio/gru/service"
)

type updatePackage struct {
	service    string
	metricType enum.MetricType
	metric     string
	toAdd      []float64
}

var (
	servicesMetrics map[string]ServiceMetric
	ch_update       chan updatePackage
	ch_compute      chan struct{}
	ch_avg          chan map[string]data.MetricStats
)

func init() {
	ch_update = make(chan updatePackage, 100)
	ch_compute = make(chan struct{})
	ch_avg = make(chan map[string]data.MetricStats)
}

func Initialize() {
	servicesList := srv.List()
	servicesMetrics = make(map[string]ServiceMetric, len(servicesList))
	for _, service := range servicesList {
		metric := ServiceMetric{
			BaseMetric: make(map[string][]float64),
			UserMetric: make(map[string][]float64),
		}

		servicesMetrics[service] = metric
	}
}

func UpdateCpuMetric(service string, toAdd []float64) {
	ch_update <- updatePackage{
		service:    service,
		metricType: enum.METRIC_T_BASE,
		metric:     enum.METRIC_CPU.ToString(),
		toAdd:      toAdd,
	}
}

func UpdateMemMetric(service string, toAdd []float64) {
	ch_update <- updatePackage{
		service:    service,
		metricType: enum.METRIC_T_BASE,
		metric:     enum.METRIC_MEM.ToString(),
		toAdd:      toAdd,
	}
}

func UpdateUserMetric(service string, metric string, toAdd []float64) {
	ch_update <- updatePackage{
		service:    service,
		metricType: enum.METRIC_T_USER,
		metric:     metric,
		toAdd:      toAdd,
	}
}

func GetMetricsAvg() map[string]data.MetricStats {
	ch_compute <- struct{}{}
	avg := <-ch_avg

	return avg
}

func StartMetricCollector() {
	go metricCollector()
}

func metricCollector() {
	for {
		select {
		case update := <-ch_update:
			updateMetric(update.service, update.metricType, update.metric, update.toAdd)
		case <-ch_compute:
			avg := computeMetrics()
			ch_avg <- avg
			clearMetrics()
		}
	}
}

func updateMetric(service string, metricType enum.MetricType, metric string, toAdd []float64) {
	if toUpdateService, ok := servicesMetrics[service]; ok {
		switch metricType {
		case enum.METRIC_T_BASE:
			values := toUpdateService.BaseMetric[metric]
			values = append(values, toAdd...)
			toUpdateService.BaseMetric[metric] = values
		case enum.METRIC_T_USER:
			if values, ok := toUpdateService.UserMetric[metric]; ok {
				values = append(values, toAdd...)
				toUpdateService.UserMetric[metric] = values
			} else {
				log.WithFields(log.Fields{
					"service": service,
					"metric":  metric,
				}).Errorln("Cannot update service metric: unknown metric")
			}
		}

	} else {
		log.WithField("service", service).Errorln("Cannot update service metric: unknown service")
	}
}

func computeMetrics() map[string]data.MetricStats {

	servicesAvg := make(map[string]data.MetricStats, len(servicesMetrics))

	for service, metrics := range servicesMetrics {
		baseAvg := make(map[string]float64, len(metrics.BaseMetric))
		for metric, values := range metrics.BaseMetric {
			value := mean(values)
			baseAvg[metric] = value
		}

		userAvg := make(map[string]float64, len(metrics.UserMetric))
		for metric, values := range metrics.UserMetric {
			value := mean(values)
			userAvg[metric] = value
		}

		serviceAvg := data.MetricStats{
			BaseMetrics: baseAvg,
			UserMetrics: userAvg,
		}

		servicesAvg[service] = serviceAvg
	}

	return servicesAvg
}

func mean(values []float64) float64 {
	if len(values) < 1 {
		return 0.0
	}

	sum := 0.0
	for _, value := range values {
		sum += value
	}

	return (sum / float64(len(values)))
}

func clearMetrics() {
	capacity := len(servicesMetrics)
	servicesMetrics = make(map[string]ServiceMetric, capacity)
}
