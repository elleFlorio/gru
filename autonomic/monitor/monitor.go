package monitor

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/autonomic/monitor/logreader"
	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/node"
	"github.com/elleFlorio/gru/service"
	"github.com/elleFlorio/gru/storage"
	"github.com/elleFlorio/gru/utils"
)

var (
	monitorActive  bool
	gruStats       GruStats
	history        statsHistory
	c_stop         chan struct{}
	c_err          chan error
	ErrNoIndexById error = errors.New("No index for such Id")
)

//History window
const W_SIZE = 50
const W_MULT = 1000

func init() {
	gruStats = GruStats{
		Service:  make(map[string]ServiceStats),
		Instance: make(map[string]InstanceStats),
	}
	history = statsHistory{make(map[string]instanceHistory)}
}

func Run() GruStats {
	log.WithField("status", "init").Debugln("Gru Monitor")
	defer log.WithField("status", "done").Debugln("Gru Monitor")
	snapshot := GruStats{
		Service:  make(map[string]ServiceStats),
		Instance: make(map[string]InstanceStats),
	}

	computeServicesStats(&gruStats)
	computeSystemCpu(&gruStats)
	makeSnapshot(&gruStats, &snapshot)
	updateServicesInstances(snapshot)
	err := saveStats(snapshot)
	if err != nil {
		log.WithField("err", "Stats data not saved").Errorln("Running monitor")
	}

	services := service.List()
	for _, name := range services {
		resetEventsStats(name, &gruStats)
		resetMetricStats(name, &gruStats)
	}

	displayStatsOfServices(snapshot)
	return snapshot
}

func computeServicesStats(stats *GruStats) {
	metrics := metric.Manager().GetMetrics()
	for name, _ := range gruStats.Service {
		updateRunningInstances(name, &gruStats, 25)
		computeServiceCpuPerc(name, &gruStats)
		computeServiceMemory(name, &gruStats)
		updateServiceMetrics(name, metrics[name], &gruStats)
	}
}

//FIXME need to check if all the windows are actually ready
func updateRunningInstances(name string, stats *GruStats, wsize int) {
	srvStats := stats.Service[name]
	pending := srvStats.Instances.Pending
	for _, inst := range pending {
		if len(history.instance[inst].cpu.sysUsage.Slice()) >= wsize {
			addResource(inst, name, "running", stats, &history)

			log.WithFields(log.Fields{
				"service": name,
				"id":      inst,
			}).Debugln("Promoted resource to running state")
		}
	}
}

func computeServiceCpuPerc(name string, stats *GruStats) {
	sum := 0.0
	avg := 0.0
	srvStats := stats.Service[name]

	if len(srvStats.Instances.Running) > 0 {
		for _, id := range srvStats.Instances.Running {
			instCpus := history.instance[id].cpu.totalUsage.Slice()
			sysCpus := history.instance[id].cpu.sysUsage.Slice()
			instCpuAvg := computeInstanceCpuPerc(instCpus, sysCpus)
			inst := stats.Instance[id]
			inst.Cpu = instCpuAvg
			stats.Instance[id] = inst
			sum += instCpuAvg

			log.WithFields(log.Fields{
				"instance": id,
				"CPUavg":   instCpuAvg,
			}).Debugln("Computed local instance average CPU")
		}

		avg = sum / float64(len(srvStats.Instances.Running))
	}

	srvStats.Cpu.Avg = avg
	srvStats.Cpu.Tot = sum
	stats.Service[name] = srvStats

	log.WithFields(log.Fields{
		"Service": name,
		"CPUavg":  avg,
		"CPUsum":  sum,
	}).Debugln("Computed local service CPU usage")
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
			cpu = (instDelta / sysDelta) * float64(node.Config().Resources.TotalCpus)
		}
		sum += cpu
	}
	return math.Min(1.0, sum/float64(len(instCpus)-1))
}

func computeServiceMemory(name string, stats *GruStats) {
	sum := 0.0
	avg := 0.0
	srv, _ := service.GetServiceByName(name)
	memLimit := srv.Configuration.Memory
	srvStats := stats.Service[name]

	if len(srvStats.Instances.Running) > 0 {
		for _, id := range srvStats.Instances.Running {
			instMem := history.instance[id].mem.Slice()
			instMemPerc := computeInstaceMemPerc(instMem, memLimit)
			inst := stats.Instance[id]
			inst.Memory = instMemPerc
			stats.Instance[id] = inst
			sum += instMemPerc

			log.WithFields(log.Fields{
				"instance": id,
				"Memory":   instMemPerc,
			}).Debugln("Computed local instance Memory")
		}

		avg = sum / float64(len(srvStats.Instances.Running))
	}

	srvStats.Memory.Avg = avg
	// If the service does not have a memory limit,
	// the sum of the memory used by all the instances
	// can exceed 100%. In this case the total memory used
	// is limited by the system memory and the total is
	// virtually equal to the average
	if memLimit != "" {
		srvStats.Memory.Tot = sum
	} else {
		srvStats.Memory.Tot = avg
	}
	stats.Service[name] = srvStats

	log.WithFields(log.Fields{
		"Service":   name,
		"MemoryAvg": avg,
		"MemorySum": sum,
	}).Debugln("Computed local service memory usage")
}

func computeInstaceMemPerc(instMem []float64, limit string) float64 {
	var err error
	totalMemory := node.Config().Resources.TotalMemory
	sum := 0.0
	avg := 0.0
	limitBytes := totalMemory
	if limit != "" {
		limitBytes, err = utils.RAMInBytes(limit)
		if err != nil {
			log.WithField("err", err).Errorln("Error computing instance memory usage")
			limitBytes = totalMemory
		}
	}

	for _, m := range instMem {
		sum += m
	}
	avg = sum / float64(len(instMem))

	return avg / float64(limitBytes)
}

func updateServiceMetrics(name string, metrics map[string][]float64, stats *GruStats) {
	if len(metrics) == 0 {
		log.Debugln("No metrics to update for service ", name)
		return
	}

	srv := stats.Service[name]
	for metric, value := range metrics {
		switch metric {
		case "execution_time":
			srv.Metrics.ResponseTime = value
			log.WithField("execution_time", value).Debugln("Updated execution time of service ", name)
		case "response_time":
			//TODO
		default:
			log.WithField("metric", metric).Errorln("Cannot update undefined metric of service ", name)
		}
	}
	stats.Service[name] = srv
}

func computeSystemCpu(stats *GruStats) {
	sum := 0.0
	for _, value := range stats.Service {
		sum += value.Cpu.Tot
	}
	stats.System.Cpu = sum
}

//TODO maybe I can just compute historical data without make a deep copy
// since now I'm serializing the structure in a string of bytes...
// check this possibility.
func makeSnapshot(src *GruStats, dst *GruStats) {
	// Copy service stats
	for name, stats := range src.Service {
		srv_src := stats
		// Copy instances status
		status_src := stats.Instances
		all_dst := make([]string, len(status_src.All), len(status_src.All))
		runnig_dst := make([]string, len(status_src.Running), len(status_src.Running))
		pending_dst := make([]string, len(status_src.Pending), len(status_src.Pending))
		stopped_dst := make([]string, len(status_src.Stopped), len(status_src.Stopped))
		paused_dst := make([]string, len(status_src.Paused), len(status_src.Paused))
		copy(all_dst, status_src.All)
		copy(runnig_dst, status_src.Running)
		copy(pending_dst, status_src.Pending)
		copy(stopped_dst, status_src.Stopped)
		copy(paused_dst, status_src.Paused)
		status_dst := service.InstanceStatus{
			all_dst,
			runnig_dst,
			pending_dst,
			stopped_dst,
			paused_dst,
		}
		// Copy events (NEEDED?)
		events_src := srv_src.Events
		stop_dst := make([]string, len(events_src.Stop), len(events_src.Stop))
		start_dst := make([]string, len(events_src.Start), len(events_src.Start))
		copy(start_dst, events_src.Start)
		copy(stop_dst, events_src.Stop)
		events_dst := EventStats{
			start_dst,
			stop_dst,
		}

		// Copy cpu
		cpu_dst := CpuStats{stats.Cpu.Avg, stats.Cpu.Tot}

		// Copy memory
		mem_dst := MemoryStats{stats.Memory.Avg, stats.Memory.Tot}

		// Copy metrics
		metrics_src := srv_src.Metrics
		respTime_dst := make([]float64, len(metrics_src.ResponseTime), len(metrics_src.ResponseTime))
		copy(respTime_dst, metrics_src.ResponseTime)
		metrics_dst := MetricStats{respTime_dst}

		srv_dst := ServiceStats{status_dst, events_dst, cpu_dst, mem_dst, metrics_dst}
		dst.Service[name] = srv_dst
	}

	//Copy instance stats
	for id, value := range src.Instance {
		inst_dst := InstanceStats{value.Cpu, value.Memory}
		dst.Instance[id] = inst_dst
	}

	//Copy system stats
	sys_status_src := src.System.Instances
	sys_all_dst := make([]string, len(sys_status_src.All), len(sys_status_src.All))
	sys_runnig_dst := make([]string, len(sys_status_src.Running), len(sys_status_src.Running))
	sys_pending_dst := make([]string, len(sys_status_src.Pending), len(sys_status_src.Pending))
	sys_stopped_dst := make([]string, len(sys_status_src.Stopped), len(sys_status_src.Stopped))
	sys_paused_dst := make([]string, len(sys_status_src.Paused), len(sys_status_src.Paused))
	copy(sys_all_dst, sys_status_src.All)
	copy(sys_runnig_dst, sys_status_src.Running)
	copy(sys_pending_dst, sys_status_src.Pending)
	copy(sys_stopped_dst, sys_status_src.Stopped)
	copy(sys_paused_dst, sys_status_src.Paused)
	sys_status_dst := service.InstanceStatus{
		sys_all_dst,
		sys_runnig_dst,
		sys_pending_dst,
		sys_stopped_dst,
		sys_paused_dst,
	}
	dst.System.Instances = sys_status_dst
	dst.System.Cpu = src.System.Cpu
}

func updateServicesInstances(stats GruStats) {
	for name, item := range stats.Service {
		srv, _ := service.GetServiceByName(name)
		srv.Instances = item.Instances
	}
}

func saveStats(stats GruStats) error {
	data, err := convertStatsToData(stats)
	if err != nil {
		log.WithField("err", err).Debugln("Cannot convert stats to data")
		return err
	} else {
		storage.StoreLocalData(data, enum.STATS)
	}

	return nil
}

func convertStatsToData(stats GruStats) ([]byte, error) {
	data, err := json.Marshal(stats)
	if err != nil {
		log.WithField("err", err).Errorln("Error marshaling stats data")
		return nil, err
	}

	return data, nil
}

func resetEventsStats(srvName string, stats *GruStats) {
	srvStats := stats.Service[srvName]

	log.WithFields(log.Fields{
		"service": srvName,
		"start":   srvStats.Events.Start,
		"stop":    srvStats.Events.Stop,
	}).Debugln("Monitored events")

	srvStats.Events = EventStats{}
	stats.Service[srvName] = srvStats
}

func resetMetricStats(srvName string, stats *GruStats) {
	srvStats := stats.Service[srvName]
	srvStats.Metrics = MetricStats{}
	stats.Service[srvName] = srvStats
}

func displayStatsOfServices(stats GruStats) {
	for srv, value := range stats.Service {
		log.WithFields(log.Fields{
			"pending:": len(value.Instances.Pending),
			"running:": len(value.Instances.Running),
			"stopped:": len(value.Instances.Stopped),
			"paused:":  len(value.Instances.Paused),
			"cpu avg":  fmt.Sprintf("%.2f", value.Cpu.Avg),
			"cpu tot":  fmt.Sprintf("%.2f", value.Cpu.Tot),
			"mem avg":  fmt.Sprintf("%.2f", value.Memory.Avg),
			"mem tot":  fmt.Sprintf("%.2f", value.Memory.Tot),
		}).Infoln("Stats computed: ", srv)
	}
}

func GetMonitorData() (GruStats, error) {
	stats := GruStats{}
	dataStats, err := storage.GetLocalData(enum.STATS)
	if err != nil {
		log.WithField("err", err).Errorln("Cannot retrieve stats data.")
	} else {
		stats, err = convertDataToStats(dataStats)
	}

	return stats, err
}

func convertDataToStats(data []byte) (GruStats, error) {
	stats := GruStats{}
	err := json.Unmarshal(data, &stats)
	if err != nil {
		log.WithField("err", err).Errorln("Error unmarshaling stats data")
	}

	return stats, err
}
