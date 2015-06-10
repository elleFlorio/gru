package analyzer

import (
	log "github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/autonomic/monitor"
	"github.com/elleFlorio/gru/service"
)

type analyzer struct{ c_err chan error }

var gruAnalytics GruAnalytics

func init() {
	resetAnalytics()
}

func resetAnalytics() {
	gruAnalytics = GruAnalytics{
		Service:  make(map[string]ServiceAnalytics),
		Instance: make(map[string]InstanceAnalytics),
	}
}

func NewAnalyzer(c_err chan error) *analyzer {
	return &analyzer{c_err}
}

func (p *analyzer) Run(stats monitor.GruStats) GruAnalytics {
	log.WithField("status", "start").Debugln("Running analyzer")
	defer log.WithField("status", "done").Debugln("Running analyzer")

	for _, name := range service.List() {
		if gruAnalytics.System.Cpu > 0 {
			computeCpuAvg(name, &stats)

			log.WithFields(log.Fields{
				"status":  "analyzing",
				"service": name,
				"CpuAvg":  gruAnalytics.Service[name].CpuAvg,
			}).Debugln("Running analyzer")
		}

		updateAnalytics(name, &stats)
	}

	gruAnalytics.System.Cpu = stats.System.Cpu

	return gruAnalytics
}

func updateAnalytics(name string, stats *monitor.GruStats) {

	srvAnalytics := gruAnalytics.Service[name]
	srvAnalytics.Instances = stats.Service[name].Instances
	gruAnalytics.Service[name] = srvAnalytics

	toBeRemoved := stats.Service[name].Events.Die

	for _, died := range toBeRemoved {
		delete(gruAnalytics.Instance, died)
	}

	for _, id := range srvAnalytics.Instances {
		instAnalytics := gruAnalytics.Instance[id]
		instAnalytics.Cpu = stats.Instance[id].Cpu
		gruAnalytics.Instance[id] = instAnalytics
	}
}

func computeCpuAvg(name string, stats *monitor.GruStats) {
	sum := 0.0
	sysOld := gruAnalytics.System.Cpu
	sysNew := stats.System.Cpu

	srvAnalytics := gruAnalytics.Service[name]

	instances := srvAnalytics.Instances
	for _, id := range instances {
		instAnalytics := gruAnalytics.Instance[id]
		instOld := instAnalytics.Cpu
		instNew := stats.Instance[id].Cpu
		// 100 * ?
		instAnalytics.CpuPerc = float64(instNew-instOld) / float64(sysNew-sysOld)
		gruAnalytics.Instance[id] = instAnalytics
		sum += instAnalytics.CpuPerc
	}

	srvAnalytics.CpuAvg = sum / float64(len(instances))
	gruAnalytics.Service[name] = srvAnalytics
}
