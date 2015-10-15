package analyzer

import (
	"errors"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/autonomic/monitor"
	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/node"
	"github.com/elleFlorio/gru/service"
	"github.com/elleFlorio/gru/storage"
	"github.com/elleFlorio/gru/utils"
)

const dataType string = "analytics"
const statsDataType string = "stats"
const overcommitratio float64 = 0.25

var (
	gruAnalytics          GruAnalytics
	ErrNoRunningInstances error = errors.New("No active instance to analyze")
)

func init() {
	gruAnalytics = GruAnalytics{
		Service: make(map[string]ServiceAnalytics),
	}
}

func Run() {
	log.WithField("status", "start").Infoln("Running analyzer")
	defer log.WithField("status", "done").Infoln("Running analyzer")

	stats, err := retrieveStats()
	if err != nil {
		log.WithField("error", "Cannot compute analytics").Errorln("Running Analyzer.")
	} else {
		updateNodeResources()
		analyzeServices(&gruAnalytics, stats)
		analyzeSystem(&gruAnalytics, stats)
		computeNodeHealth(&gruAnalytics)
		analyzeCluster(&gruAnalytics)
		err = saveAnalytics(gruAnalytics)
		if err != nil {
			log.WithField("error", "Cluster analytics data not saved ").Errorln("Running Analyzer")
		}
	}
}

func retrieveStats() (monitor.GruStats, error) {
	stats := monitor.GruStats{}
	dataStats, err := storage.GetLocalData(enum.STATS)
	if err != nil {
		log.WithField("error", err).Errorln("Cannot retrieve stats data.")
	} else {
		stats, err = monitor.ConvertDataToStats(dataStats)
	}

	return stats, err
}

func updateNodeResources() {
	_, err := node.UsedCpus()
	if err != nil {
		log.WithField("error", err).Errorln("Computing node used CPU")
	}
	_, err = node.UsedMemory()
	if err != nil {
		log.WithField("error", err).Errorln("Computing node used memory")
	}
}

func analyzeServices(analytics *GruAnalytics, stats monitor.GruStats) {
	for name, value := range stats.Service {
		//TODO some labels are not yet ready
		temp := 0.0
		loadLabel := enum.FromValue(temp)
		cpuLabel := enum.FromValue(value.Cpu.Tot)
		memLabel := enum.FromValue(temp)
		srvResources := computeServiceResources(name)
		instances := value.Instances

		health := (loadLabel + cpuLabel + memLabel + srvResources) / 4 //I don't like this...

		srvRes := ResourcesAnalytics{
			cpuLabel,
			memLabel,
			srvResources,
		}

		srvAnalytics := ServiceAnalytics{
			loadLabel,
			srvRes,
			instances,
			health,
		}

		analytics.Service[name] = srvAnalytics
	}
}

func computeServiceResources(name string) enum.Label {
	nodeMem := node.Config().Resources.TotalMemory
	nodeCpu := node.Config().Resources.TotalCpus
	nodeUsedMem := node.Config().Resources.UsedMemory
	nodeUsedCpu := node.Config().Resources.UsedCpu

	srv, _ := service.GetServiceByName(name)
	//I use CpuSetValue despite Swarm
	//srvCpu := srv.Configuration.CpuShares
	srvCpu := srv.Configuration.CpuSet
	srvMem, err := utils.RAMInBytes(srv.Configuration.Memory)
	if err != nil {
		log.WithField("error", err).Warnln("Cannot convert service RAM in Bytes.")
		return enum.RED
	}

	var (
		cpuScore float64 = 1
		memScore float64 = 1
		weight   float64 = 1
	)

	if nodeMem < int64(srvMem) || nodeCpu < int64(srvCpu) {
		return enum.RED
	}

	nodeCpuOverCommit := (float64(nodeCpu) * overcommitratio) + float64(nodeCpu)
	nodeMemOverCommit := (float64(nodeMem) * overcommitratio) + float64(nodeMem)

	if srvCpu > 0 {
		cpuScore = float64(nodeUsedCpu+srvCpu) / nodeCpuOverCommit
	}

	if srvMem > 0 {
		memScore = float64(nodeUsedMem+srvMem) / nodeMemOverCommit
	}

	if cpuScore <= 1.0 && memScore <= 1.0 {
		weight = (cpuScore + memScore) / 2
	}

	return enum.FromValue(weight)

}

func analyzeSystem(analytics *GruAnalytics, stats monitor.GruStats) {
	sysSrvs := []string{}
	for name, _ := range stats.Service {
		sysSrvs = append(sysSrvs, name)
	}

	temp := 0.0
	cpuLabel := enum.FromValue(stats.System.Cpu)
	memLabel := enum.FromValue(temp)
	resLabel := computeSystemResources()
	instances := stats.System.Instances

	health := resLabel //Ok, maybe this is a bit... "mah"...

	sysRes := ResourcesAnalytics{
		cpuLabel,
		memLabel,
		resLabel,
	}

	SystemAnalytics := SystemAnalytics{
		sysSrvs,
		sysRes,
		instances,
		health,
	}

	gruAnalytics.System = SystemAnalytics
}

func computeSystemResources() enum.Label {
	totalCpu := float64(node.Config().Resources.TotalCpus)
	totalMemory := float64(node.Config().Resources.TotalMemory)
	usedCpu := float64(node.Config().Resources.UsedCpu)
	usedMemory := float64(node.Config().Resources.UsedMemory)

	cpuRatio := usedCpu / totalCpu
	memRatio := usedMemory / totalMemory
	avgRatio := (cpuRatio + memRatio) / 2

	return enum.FromValue(avgRatio)
}

func computeNodeHealth(analytics *GruAnalytics) {
	nServices := len(analytics.Service)
	sumHealth := 0.0
	for _, value := range analytics.Service {
		sumHealth += value.Health.Value()
	}
	srvAvgHealth := sumHealth / float64(nServices)

	sysHealth := analytics.System.Health.Value()

	totHealth := (srvAvgHealth + sysHealth) / 2
	totHealthLabel := enum.FromValue(totHealth)

	analytics.Health = totHealthLabel
}

func analyzeCluster(analytics *GruAnalytics) {
	peers := getPeersAnalytics()
	computeServicesAvg(peers, analytics)
	computeClusterAvg(peers, analytics)
}

func getPeersAnalytics() []GruAnalytics {
	peers := make([]GruAnalytics, 0)
	dataAn, _ := storage.GetAllData(enum.ANALYTICS)
	for _, data := range dataAn {
		a, _ := ConvertDataToAnalytics(data)
		peers = append(peers, a)
	}

	return peers
}

func computeServicesAvg(peers []GruAnalytics, analytics *GruAnalytics) {
	avg := make(map[string]ServiceAnalytics)

	for _, name := range service.List() {
		active := []ServiceAnalytics{}
		var avgSa ServiceAnalytics

		if ls, ok := analytics.Service[name]; ok {
			active = append(active, ls)
		}

		for _, peer := range peers {
			if ps, ok := peer.Service[name]; ok {
				active = append(active, ps)
			}
		}

		if len(active) > 1 {
			avgSa = active[0]
			active = active[1:]

			sumLoad := avgSa.Load.Value()
			sumCpu := avgSa.Resources.Cpu.Value()
			sumMem := avgSa.Resources.Memory.Value()
			sumH := avgSa.Health.Value()

			for _, actv := range active {
				//LABELS
				sumLoad += actv.Load.Value()
				sumCpu += actv.Resources.Cpu.Value()
				sumMem += actv.Resources.Memory.Value()
				sumH += actv.Health.Value()

				//INSTANCES
				avgSa.Instances.All = append(avgSa.Instances.All, actv.Instances.All...)
				avgSa.Instances.Running = append(avgSa.Instances.Running, actv.Instances.Running...)
				avgSa.Instances.Pending = append(avgSa.Instances.Pending, actv.Instances.Pending...)
				avgSa.Instances.Stopped = append(avgSa.Instances.Stopped, actv.Instances.Stopped...)
				avgSa.Instances.Paused = append(avgSa.Instances.Paused, actv.Instances.Paused...)
			}

			total := float64(len(active) + 1)
			avgLoad := sumLoad / total
			avgCpu := sumCpu / total
			avgMem := sumMem / total
			avgH := sumH / total

			avgSa.Load = enum.FromLabelValue(avgLoad)
			avgSa.Resources.Cpu = enum.FromLabelValue(avgCpu)
			avgSa.Resources.Memory = enum.FromLabelValue(avgMem)
			avgSa.Health = enum.FromLabelValue(avgH)

			avg[name] = avgSa

		} else if len(active) == 1 {
			avg[name] = active[0]
		}
	}

	analytics.Service = avg
}

func computeClusterAvg(peers []GruAnalytics, analytics *GruAnalytics) {
	clstrSrvs := []string{}
	var sumCpu float64 = 0
	var sumMem float64 = 0
	var sumH float64 = 0

	for _, peer := range peers {
		clstrSrvs = checkAndAppend(clstrSrvs, peer.System.Services)
		sumCpu += peer.System.Resources.Cpu.Value()
		sumMem += peer.System.Resources.Memory.Value()
		sumH += peer.System.Health.Value()
	}

	clstrSrvs = checkAndAppend(clstrSrvs, analytics.System.Services)
	total := float64(len(peers) + 1)
	avgCpu := (analytics.System.Resources.Cpu.Value() + sumCpu) / total
	avgMem := (analytics.System.Resources.Memory.Value() + sumMem) / total
	avgH := (analytics.System.Health.Value() + sumH) / total

	analytics.Cluster.Services = clstrSrvs
	analytics.Cluster.ResourcesAnalytics.Cpu = enum.FromLabelValue(avgCpu)
	analytics.Cluster.ResourcesAnalytics.Memory = enum.FromLabelValue(avgMem)
	analytics.Cluster.Health = enum.FromLabelValue(avgH)
}

func checkAndAppend(slice []string, values []string) []string {
	var notContains bool
	for _, value := range values {
		notContains = true
		for _, item := range slice {
			if item == value {
				notContains = false
			}
		}

		if notContains {
			slice = append(slice, value)
		}
	}

	return slice
}

// This is trivial, but improve readability
func saveAnalytics(analytics GruAnalytics) error {
	data, err := convertAnalyticsToData(analytics)
	if err != nil {
		log.WithField("error", "Cannot convert analytics to data").Debugln("Running Analyzer")
		return err
	} else {
		storage.StoreData(enum.CLUSTER.ToString(), data, enum.ANALYTICS)
	}

	return nil
}
