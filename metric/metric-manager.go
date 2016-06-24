package metric

import (
	"errors"
	"time"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/data"
	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/service"
)

type metricService interface {
	Name() string
	Initialize(map[string]interface{}) error
	StoreMetrics(GruMetric) error
}

const c_AUTO_LOOP_EPSILON = 5

var (
	metricServices  []metricService
	metricServ      int
	metrics         GruMetric
	ErrNotSupported error = errors.New("Metric service not supported")
)

func init() {
	metricServices = []metricService{
		&noService{},
		&influxdb{},
	}
}

func New(name string, conf map[string]interface{}) (metricService, error) {
	metricServ = 0
	for index, mtrc := range metricServices {
		if mtrc.Name() == name {
			err := mtrc.Initialize(conf)
			if err != nil {
				log.WithFields(log.Fields{
					"err":     err,
					"service": mtrc.Name(),
				}).Errorln("Error initializing metric service")
				return metricServices[metricServ], err
			}
			metricServ = index
			log.WithField("service", name).Debugln("Initialized metric service")
			return metricServices[metricServ], nil
		}
	}

	return metricServices[metricServ], ErrNotSupported
}

func Name() string {
	return activeService().Name()
}

func Initialize(conf map[string]interface{}) error {
	return activeService().Initialize(conf)
}

func activeService() metricService {
	return metricServices[metricServ]
}

func Metrics() GruMetric {
	return metrics
}

func StartMetricCollector() {
	go startCollector()
}

func startCollector() {
	interval := cfg.GetAgentAutonomic().LoopTimeInterval + c_AUTO_LOOP_EPSILON
	ticker := time.NewTicker(time.Duration(interval) * time.Second)

	for {
		select {
		case <-ticker.C:
			collectMetrics()
		}
	}
}

func collectMetrics() {
	log.Debugln("Collecting metrics")
	updateMetrics()
	err := storeMetrics(metrics)
	if err != nil {
		log.WithField("errr", err).Errorln("Error collecting agent metrics")
	}
}

func updateMetrics() {
	var err error
	metrics = newMetrics()
	metrics.Node.UUID = cfg.GetNodeConfig().UUID
	metrics.Node.Name = cfg.GetNodeConfig().Name

	for _, name := range service.List() {
		srv, _ := service.GetServiceByName(name)
		srv_metrics := ServiceMetric{}
		srv_metrics.Name = name
		srv_metrics.Image = srv.Image
		srv_metrics.Type = srv.Type

		srv_metrics.Instances.All = len(srv.Instances.All)
		srv_metrics.Instances.Pending = len(srv.Instances.Pending)
		srv_metrics.Instances.Running = len(srv.Instances.Running)
		srv_metrics.Instances.Paused = len(srv.Instances.Paused)
		srv_metrics.Instances.Stopped = len(srv.Instances.Stopped)

		stats, err := data.GetStats()
		if err != nil {
			log.WithField("err", err).Warnln("Cannot update stats metrics")
		} else {
			if srv_stats, ok := stats.Metrics.Service[name]; ok {
				srv_metrics.Stats = srv_stats
			} else {
				log.Warnln("Cannot find stats metrics for service ", name)
			}

			metrics.Node.Stats = stats.Metrics.System
		}

		analytics, err := data.GetAnalytics()
		if err != nil {
			log.WithField("err", err).Warnln("Cannot update analytics metrics")
		} else {
			if srv_analytisc, ok := analytics.Service[name]; ok {
				srv_metrics.Analytics = srv_analytisc
			} else {
				log.Debugln("Cannot find analytics metrics for service ", name)
			}
		}

		shared, err := data.GetSharedCluster()
		if err != nil {
			log.WithField("err", err).Warnln("Cannot update shared data metrics")
		} else {
			if srv_shared, ok := shared.Service[name]; ok {
				srv_metrics.Shared = srv_shared.Data
			}
		}

		metrics.Service[name] = srv_metrics
	}

	plc, err := data.GetPolicy()
	if err != nil {
		log.WithField("err", err).Warnln("Cannot update plans metrics")
	} else {
		metrics.Policy.Name = plc.Name
		metrics.Policy.Weight = plc.Weight
	}

}

func newMetrics() GruMetric {
	metricsEmpty := GruMetric{Service: make(map[string]ServiceMetric)}
	metricDataEmpty := data.MetricData{
		BaseMetrics: map[string]float64{
			enum.METRIC_CPU_AVG.ToString(): 0.0,
			enum.METRIC_MEM_AVG.ToString(): 0.0,
		},
		UserMetrics: map[string]float64{
			"nodata": 0.0,
		},
	}
	analyticsDataEmpty := data.AnalyticData{
		BaseAnalytics: map[string]float64{
			enum.METRIC_CPU_AVG.ToString(): 0.0,
			enum.METRIC_MEM_AVG.ToString(): 0.0,
		},
		UserAnalytics: map[string]float64{
			"nodata": 0.0,
		},
	}
	sharedDataEmpty := data.SharedData{
		BaseShared: map[string]float64{
			enum.METRIC_CPU_AVG.ToString(): 0.0,
			enum.METRIC_MEM_AVG.ToString(): 0.0,
		},
		UserShared: map[string]float64{
			"nodata": 0.0,
		},
	}

	metricsEmpty.Node = NodeMetrics{
		Stats: metricDataEmpty,
	}

	for _, name := range service.List() {
		metricsEmpty.Service[name] = ServiceMetric{
			Stats:     metricDataEmpty,
			Analytics: analyticsDataEmpty,
			Shared:    sharedDataEmpty,
		}
	}

	policyEmpty := PolicyMetric{Name: "noaction", Weight: 1.0}
	metricsEmpty.Policy = policyEmpty

	return metricsEmpty
}

func storeMetrics(metrics GruMetric) error {
	return activeService().StoreMetrics(metrics)
}
