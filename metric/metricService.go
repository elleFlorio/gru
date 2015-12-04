package metric

import (
	"errors"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/autonomic/analyzer"
	"github.com/elleFlorio/gru/autonomic/monitor"
	"github.com/elleFlorio/gru/autonomic/planner"
	"github.com/elleFlorio/gru/node"
	"github.com/elleFlorio/gru/service"
)

type metricService interface {
	Name() string
	Initialize(map[string]interface{}) error
	StoreMetrics(GruMetric) error
}

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

func clearMetrics() {
	metrics = GruMetric{Service: make(map[string]ServiceMetric)}
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

func StoreMetrics(metrics GruMetric) error {
	return activeService().StoreMetrics(metrics)
}

func activeService() metricService {
	return metricServices[metricServ]
}

func Metrics() GruMetric {
	return metrics
}

func UpdateMetrics() {
	var err error
	clearMetrics()
	metrics.Node.UUID = node.Config().UUID

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

		stats, err := monitor.GetMonitorData()
		if err != nil {
			log.WithField("err", err).Errorln("Cannot update stats metrics")
		} else {
			if srv_stats, ok := stats.Service[name]; ok {
				srv_metrics.Stats.CpuAvg = srv_stats.Cpu.Avg
				srv_metrics.Stats.CpuTot = srv_stats.Cpu.Tot
				srv_metrics.Stats.MemAvg = srv_stats.Memory.Avg
				srv_metrics.Stats.MemTot = srv_stats.Memory.Tot

				metrics.Node.Cpu = stats.System.Cpu
				metrics.Node.Memory = 0.0 // TODO
			} else {
				log.Warnln("Cannot find stats metrics for service ", name)
			}
		}

		analytics, err := analyzer.GetAnalyzerData()
		if err != nil {
			log.WithField("err", err).Errorln("Cannot update analytics metrics")
		} else {
			if srv_analytisc, ok := analytics.Service[name]; ok {
				srv_metrics.Analytics.Cpu = srv_analytisc.Resources.Cpu
				srv_metrics.Analytics.Memory = srv_analytisc.Resources.Memory
				srv_metrics.Analytics.Resources = srv_analytisc.Resources.Available
				srv_metrics.Analytics.Load = srv_analytisc.Load
				srv_metrics.Analytics.Health = srv_analytisc.Health
			} else {
				log.Warnln("Cannot find analytics metrics for service ", name)
			}
		}

		metrics.Service[name] = srv_metrics
	}

	plans, err := planner.GetPlannerData()
	if err != nil {
		log.WithField("err", err).Errorln("Cannot update plans metrics")
	} else {
		metrics.Plan.Policy = plans.Policy
		metrics.Plan.Target = plans.Target.Name
		metrics.Plan.Weight = plans.Weight
	}

}
