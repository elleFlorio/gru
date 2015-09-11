package analyzer

import (
	"errors"

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

func GetNodeAnalytics() GruAnalytics {
	services := GetServicesAanalytics()
	instances := GetInstancesAanalytics()
	system := GetSystemAnalytics()

	return GruAnalytics{
		services,
		instances,
		system,
	}

}

func GetServiceAanalytics(name string) ServiceAnalytics {
	return gruAnalytics.Service[name]
}

func GetServicesAanalytics() map[string]ServiceAnalytics {
	return gruAnalytics.Service
}

func GetInstanceAanalytics(id string) InstanceAnalytics {
	return gruAnalytics.Instance[id]
}

func GetInstancesAanalytics() map[string]InstanceAnalytics {
	return gruAnalytics.Instance
}

func GetSystemAnalytics() SystemAnalytics {
	return gruAnalytics.System
}

func (p *analyzer) Run(stats monitor.GruStats) GruAnalytics {
	log.WithField("status", "start").Debugln("Running analyzer")
	defer log.WithField("status", "done").Debugln("Running analyzer")

	sysCpuPerc := 0.0

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
			cpuTot, cpuAvg := computeServiceCpuPerc(name, &gruAnalytics, &stats)
			log.WithFields(log.Fields{
				"status":  "analyzing",
				"service": name,
				"CpuTot":  cpuTot,
				"CpuAvg":  cpuAvg,
			}).Debugln("Running analyzer")

			srv := gruAnalytics.Service[name]
			srv.CpuTot = cpuTot
			srv.CpuAvg = cpuAvg
			gruAnalytics.Service[name] = srv

			sysCpuPerc += cpuTot
		}
	}

	gruAnalytics.System.Cpu.CpuPerc = sysCpuPerc
	updateSystemInstances(&gruAnalytics)

	log.WithFields(log.Fields{
		"status":   "analyzing",
		"CpuTotal": sysCpuPerc,
	}).Debugln("Running analyzer")

	return gruAnalytics
}

func updateInstances(name string, analytics *GruAnalytics, stats *monitor.GruStats, numberOfData int) {
	srvStats := stats.Service[name]
	srvAnalytics := analytics.Service[name]

	srvAnalytics.Instances.All = srvStats.Instances.All
	active, pending := getActiveInstances(srvStats.Instances.Running, stats, numberOfData)

	srvAnalytics.Instances.Active = active
	srvAnalytics.Instances.Pending = pending
	srvAnalytics.Instances.Stopped = srvStats.Instances.Stopped
	srvAnalytics.Instances.Paused = srvStats.Instances.Paused

	log.WithFields(log.Fields{
		"status":  "instance updated",
		"service": name,
		"all":     len(srvAnalytics.Instances.All),
		"pending": len(srvAnalytics.Instances.Pending),
		"active":  len(srvAnalytics.Instances.Active),
		"stopped": len(srvAnalytics.Instances.Stopped),
		"paused":  len(srvAnalytics.Instances.Paused),
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

func computeServiceCpuPerc(name string, analytics *GruAnalytics, stats *monitor.GruStats) (float64, float64) {
	sum := 0.0
	avg := 0.0
	active := analytics.Service[name].Instances.Active

	for _, id := range active {
		instCpus := stats.Instance[id].Cpu.TotalUsage
		sysCpus := stats.Instance[id].Cpu.SysUsage
		instCpuAvg := computeInstanceCpuPerc(instCpus, sysCpus)
		inst := analytics.Instance[id]
		inst.Cpu.CpuPerc = instCpuAvg
		analytics.Instance[id] = inst
		sum += instCpuAvg
	}

	avg = sum / float64(len(active))

	return sum, avg
}

// Since linux compute the cpu usage in units of jiffies, it needs to be converted
// in % using the formula used in this function.
// Explaination: http://stackoverflow.com/questions/1420426/calculating-cpu-usage-of-a-process-in-linux
func computeInstanceCpuPerc(instCpus []float64, sysCpus []float64) float64 {
	sum := 0.0
	instNext := 0.0
	sysNext := 0.0
	instPrev := 0.0
	sysPrev := 0.0
	cpu := 0.0

	for i := 1; i < len(instCpus); i++ {
		instPrev = instCpus[i-1]
		sysPrev = sysCpus[i-1]
		instNext = instCpus[i]
		sysNext = sysCpus[i]
		instDelta := instNext - instPrev
		sysDelta := sysNext - sysPrev
		if sysDelta == 0 {
			cpu = 0
		} else {
			// "100 * cpu" should produce values in [0, 100]
			cpu = instDelta / sysDelta
		}
		sum += cpu
	}
	return sum / float64(len(instCpus)-1)
}

func updateSystemInstances(analytics *GruAnalytics) {
	analytics.System.Instances.All = make([]string, 0)
	analytics.System.Instances.Active = make([]string, 0)
	analytics.System.Instances.Pending = make([]string, 0)
	analytics.System.Instances.Paused = make([]string, 0)
	analytics.System.Instances.Stopped = make([]string, 0)

	for _, v := range analytics.Service {
		srvInst := v.Instances
		analytics.System.Instances.All = append(analytics.System.Instances.All, srvInst.All...)
		analytics.System.Instances.Active = append(analytics.System.Instances.Active, srvInst.Active...)
		analytics.System.Instances.Pending = append(analytics.System.Instances.Pending, srvInst.Pending...)
		analytics.System.Instances.Stopped = append(analytics.System.Instances.Stopped, srvInst.Stopped...)
		analytics.System.Instances.Paused = append(analytics.System.Instances.Paused, srvInst.Paused...)
	}
}
