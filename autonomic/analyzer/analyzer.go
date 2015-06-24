package analyzer

import (
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/autonomic/monitor"
	"github.com/elleFlorio/gru/service"
)

type analyzer struct {
	c_err chan error
}

var (
	gruAnalytics          GruAnalytics
	ErrNoRunningInstances error = errors.New("No active instance to analyze")
)

func init() {
	gruAnalytics = GruAnalytics{
		Service:  make(map[string]ServiceAnalytics),
		Instance: make(map[string]InstanceAnalytics),
	}
}

func NewAnalyzer(c_err chan error) *analyzer {
	return &analyzer{
		c_err,
	}
}

func (p *analyzer) Run(stats monitor.GruStats) GruAnalytics {
	log.WithField("status", "start").Debugln("Running analyzer")
	defer log.WithField("status", "done").Debugln("Running analyzer")

	for _, name := range service.List() {
		updateInstances(name, &gruAnalytics, &stats, monitor.W_SIZE)
		if len(gruAnalytics.Service[name].Instances.Active) < 1 {
			log.WithFields(log.Fields{
				"status":  "analyzing",
				"service": name,
				"pending": len(gruAnalytics.Service[name].Instances.Pending),
				"error":   ErrNoRunningInstances,
			}).Warnln("Running analyzer")
		} else {
			cpuAvg := computeServiceCpuPerc(name, &gruAnalytics, &stats)
			log.WithFields(log.Fields{
				"status":  "analyzing",
				"service": name,
				"CpuAvg":  cpuAvg,
			}).Debugln("Running analyzer")

			srv := gruAnalytics.Service[name]
			srv.CpuAvg = cpuAvg
			gruAnalytics.Service[name] = srv
		}
	}

	return gruAnalytics
}

func updateInstances(name string, analytics *GruAnalytics, stats *monitor.GruStats, numberOfData int) {
	srvStats := stats.Service[name]
	srvAnalytics := analytics.Service[name]

	srvAnalytics.Instances.All = srvStats.Instances.All
	active, pending := getActiveInstances(srvStats.Instances.Running, stats, numberOfData)
	for _, id := range srvStats.Events.Start {
		pending = append(pending, id)
	}
	srvAnalytics.Instances.Active = active
	srvAnalytics.Instances.Pending = pending
	srvAnalytics.Instances.Stopped = srvStats.Instances.Stopped
	srvAnalytics.Instances.Paused = srvStats.Instances.Paused

	log.WithFields(log.Fields{
		"status":  "instance updated",
		"service": name,
		"all":     srvAnalytics.Instances.All,
		"pending": srvAnalytics.Instances.Pending,
		"active":  srvAnalytics.Instances.Active,
		"stopped": srvAnalytics.Instances.Stopped,
		"paused":  srvAnalytics.Instances.Paused,
	}).Debugln("Running analyzer")

	analytics.Service[name] = srvAnalytics

	toBeRemoved := srvStats.Events.Stop
	for _, stopped := range toBeRemoved {
		delete(analytics.Instance, stopped)
	}
}

func getActiveInstances(running []string, stats *monitor.GruStats, numberOfData int) ([]string, []string) {
	active := []string{}
	pending := []string{}

	for _, id := range running {
		instHistory := stats.Instance[id].Cpu.TotalUsage
		if len(instHistory) < numberOfData {
			pending = append(pending, id)
		} else {
			active = append(active, id)
		}
	}

	return active, pending
}

func computeServiceCpuPerc(name string, analytics *GruAnalytics, stats *monitor.GruStats) float64 {
	sum := 0.0
	avg := 0.0
	active := analytics.Service[name].Instances.Active
	fmt.Println("Service: ", name)
	fmt.Println("Active: ", len(active))

	for _, id := range active {
		instCpus := stats.Instance[id].Cpu.TotalUsage
		sysCpus := stats.Instance[id].Cpu.SysUsage
		fmt.Println("instCpus: ", instCpus)
		fmt.Println("sysCpus: ", sysCpus)
		instCpuAvg := computeInstanceCpuPerc(instCpus, sysCpus)
		fmt.Println("instCpuAvg: ", instCpuAvg)
		inst := analytics.Instance[id]
		inst.Cpu.CpuPerc = instCpuAvg
		analytics.Instance[id] = inst
		sum += instCpuAvg
	}
	avg = sum / float64(len(active))

	return avg
}

func computeInstanceCpuPerc(instCpus []float64, sysCpus []float64) float64 {
	sum := 0.0
	instNext := 0.0
	sysNext := 0.0
	instPrev := 0.0
	sysPrev := 0.0

	for i := 1; i < len(instCpus); i++ {
		instPrev = instCpus[i-1]
		sysPrev = sysCpus[i-1]
		instNext = instCpus[i]
		sysNext = sysCpus[i]
		instDelta := instNext - instPrev
		sysDelta := sysNext - sysPrev
		cpu := 100 * instDelta / sysDelta
		sum += cpu
	}

	return sum / float64(len(instCpus)-1)
}
