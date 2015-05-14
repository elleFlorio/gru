package autonomic

import (
	log "github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/service"
)

type Analyze struct{ c_err chan error }

var gruAnalytics GruAnalytics

func init() {
	gruAnalytics = GruAnalytics{
		Service:  make(map[string]ServiceAnalytics),
		Instance: make(map[string]InstanceAnalytics),
	}
}

func (p *Analyze) run(stats GruStats) GruAnalytics {
	log.WithField("status", "start").Debugln("Running analyzer")
	defer log.WithField("status", "done").Debugln("Running analyzer")

	for _, name := range service.List() {
		if gruAnalytics.System.Cpu != 0 {
			computeCpuAvg(name, &stats)
		}
		updateAnalytics(name, &stats)
	}

	return gruAnalytics
}

func updateAnalytics(name string, stats *GruStats) {

	srvAnalytics := gruAnalytics.Service[name]
	srvAnalytics.Instances = stats.Service[name].Instances
	gruAnalytics.Service[name] = srvAnalytics

	for _, id := range srvAnalytics.Instances {
		instAnalytics := gruAnalytics.Instance[id]
		instAnalytics.Cpu = stats.Instance[id].Cpu
		gruAnalytics.Instance[id] = instAnalytics
	}

	gruAnalytics.System.Cpu = stats.System.Cpu

}

func computeCpuAvg(name string, stats *GruStats) {
	sum := 0.0
	sysOld := gruAnalytics.System.Cpu
	sysNew := stats.System.Cpu

	srvAnalytics := gruAnalytics.Service[name]

	instances := srvAnalytics.Instances
	for _, id := range instances {
		instAnalytics := gruAnalytics.Instance[id]
		instOld := instAnalytics.Cpu
		instNew := stats.Instance[id].Cpu
		instAnalytics.CpuPerc = float64(instNew-instOld) / float64(sysNew-sysOld)
		gruAnalytics.Instance[id] = instAnalytics
		sum += instAnalytics.CpuPerc
	}

	srvAnalytics.CpuAvg = sum / float64(len(instances))
	gruAnalytics.Service[name] = srvAnalytics
}
