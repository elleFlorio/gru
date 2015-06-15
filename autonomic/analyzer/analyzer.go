package analyzer

import (
	"errors"

	log "github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/autonomic/monitor"
	"github.com/elleFlorio/gru/service"
)

type analyzer struct{ c_err chan error }

var (
	gruAnalytics      GruAnalytics
	ErrNotValidCpuAvg error = errors.New("Cpu Avg is not a valid value")
)

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

		err := computeCpuAvg(name, &stats)
		if err != nil {
			log.WithFields(log.Fields{
				"status":  "analyzing",
				"error":   err,
				"service": name,
				"CpuAvg":  gruAnalytics.Service[name].CpuAvg,
			}).Warnln("Running analyzer")
		}

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
	srvStats := stats.Service[name]
	srvAnalytics := gruAnalytics.Service[name]

	srvAnalytics.Instances.All = srvStats.Instances.All
	srvAnalytics.Instances.Pending = srvStats.Events.Start
	// pending instances should not be analyzed, because we don't have
	// previous data to compare.
	srvAnalytics.Instances.Running = []string{}
	for _, running := range srvStats.Instances.Running {
		if !isPending(running, srvAnalytics.Instances.Pending) {
			srvAnalytics.Instances.Running = append(srvAnalytics.Instances.Running, running)
		}
	}
	srvAnalytics.Instances.Stopped = srvStats.Instances.Stopped
	srvAnalytics.Instances.Paused = srvStats.Instances.Paused

	log.WithFields(log.Fields{
		"status":  "instance updated",
		"service": name,
		"all":     srvAnalytics.Instances.All,
		"pending": srvAnalytics.Instances.Pending,
		"running": srvAnalytics.Instances.Running,
		"stopped": srvAnalytics.Instances.Stopped,
		"paused":  srvAnalytics.Instances.Paused,
	}).Debugln("Running analyzer")

	gruAnalytics.Service[name] = srvAnalytics

	toBeRemoved := srvStats.Events.Stop
	for _, stopped := range toBeRemoved {
		delete(gruAnalytics.Instance, stopped)
	}
}

func isPending(id string, pending []string) bool {
	for _, pndng := range pending {
		if id == pndng {
			return true
		}
	}

	return false
}

func computeCpuAvg(name string, stats *monitor.GruStats) error {
	var err error = nil
	sum := 0.0
	avg := 0.0
	sysOld := gruAnalytics.System.Cpu
	sysNew := stats.System.Cpu

	srvAnalytics := gruAnalytics.Service[name]

	instances := srvAnalytics.Instances
	nRunning := len(srvAnalytics.Instances.Running)
	for _, id := range instances.Running {
		instAnalytics := gruAnalytics.Instance[id]
		instOld := instAnalytics.Cpu
		instNew := stats.Instance[id].Cpu
		instAnalytics.CpuPerc = 100 * float64(instNew-instOld) / float64(sysNew-sysOld)
		gruAnalytics.Instance[id] = instAnalytics
		sum += instAnalytics.CpuPerc
	}
	if nRunning != 0 {
		avg = sum / float64(nRunning)
	} else {
		avg = 0
		err = ErrNotValidCpuAvg
	}

	srvAnalytics.CpuAvg = avg
	gruAnalytics.Service[name] = srvAnalytics
	return err
}

func updateAnalytics(name string, stats *monitor.GruStats) {

	for id, instSt := range stats.Instance {
		instAn := gruAnalytics.Instance[id]
		instAn.Cpu = instSt.Cpu
		gruAnalytics.Instance[id] = instAn
	}
}
