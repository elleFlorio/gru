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
		updateInstances(name, &stats)

		computeCpuAvg(name, &stats)

		log.WithFields(log.Fields{
			"status":  "analyzing",
			"service": name,
			"CpuAvg":  gruAnalytics.Service[name].CpuAvg,
		}).Debugln("Running analyzer")

		updateAnalytics(name, &stats)
	}

	gruAnalytics.System.Cpu = stats.System.Cpu

	return gruAnalytics
}

func updateInstances(name string, stats *monitor.GruStats) {
	srvAnalytics := gruAnalytics.Service[name]
	srvAnalytics.Instances.All = stats.Service[name].Instances.All
	srvAnalytics.Instances.Running = stats.Service[name].Instances.Running
	srvAnalytics.Instances.Stopped = stats.Service[name].Instances.Stopped
	srvAnalytics.Instances.Paused = stats.Service[name].Instances.Paused
	gruAnalytics.Service[name] = srvAnalytics

	toBeRemoved := stats.Service[name].Events.Stop

	for _, died := range toBeRemoved {
		delete(gruAnalytics.Instance, died)
	}
}

func computeCpuAvg(name string, stats *monitor.GruStats) {
	sum := 0.0
	counter := 0
	sysOld := gruAnalytics.System.Cpu
	sysNew := stats.System.Cpu

	srvAnalytics := gruAnalytics.Service[name]

	instances := srvAnalytics.Instances
	for _, id := range instances.Running {
		if !isJustStarted(id, stats.Service[name].Events.Start) {
			instAnalytics := gruAnalytics.Instance[id]
			instOld := instAnalytics.Cpu
			instNew := stats.Instance[id].Cpu
			instAnalytics.CpuPerc = 100 * float64(instNew-instOld) / float64(sysNew-sysOld)
			gruAnalytics.Instance[id] = instAnalytics
			sum += instAnalytics.CpuPerc
			counter++
		}
	}

	srvAnalytics.CpuAvg = sum / float64(counter)
	gruAnalytics.Service[name] = srvAnalytics
}

func isJustStarted(id string, started []string) bool {
	for _, strtd := range started {
		if id == strtd {
			return true
		}
	}

	return false
}

func updateAnalytics(name string, stats *monitor.GruStats) {

	for id, instSt := range stats.Instance {
		instAn := gruAnalytics.Instance[id]
		instAn.Cpu = instSt.Cpu
		gruAnalytics.Instance[id] = instAn
	}
}
