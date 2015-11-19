package analyzer

import (
	"encoding/json"
	"errors"
	"strings"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/autonomic/monitor"
	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/node"
	"github.com/elleFlorio/gru/service"
	"github.com/elleFlorio/gru/storage"
	"github.com/elleFlorio/gru/utils"
)

const overcommitratio float64 = 0.0

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

	stats, err := monitor.GetMonitorData()
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

		displayAnalyticsOfServices(gruAnalytics)
	}
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

	log.WithFields(log.Fields{
		"totalcpu": node.Config().Resources.TotalCpus,
		"usedcpu":  node.Config().Resources.UsedCpu,
		"totalmem": node.Config().Resources.TotalMemory,
		"usedmem":  node.Config().Resources.UsedMemory,
	}).Debugln("Updated node resources")
}

func analyzeServices(analytics *GruAnalytics, stats monitor.GruStats) {
	for name, value := range stats.Service {
		load := analyzeServiceLoad(name, value.Metrics.ResponseTime)
		cpu := value.Cpu.Tot
		mem := value.Memory.Tot
		resAvailable := resourcesAvailable(name)
		instances := value.Instances

		health := 1 - ((load + mem + cpu - resAvailable) / 4) //I don't like this...

		srvRes := ResourcesAnalytics{
			cpu,
			mem,
			resAvailable,
		}

		srvAnalytics := ServiceAnalytics{
			load,
			srvRes,
			instances,
			health,
		}

		analytics.Service[name] = srvAnalytics
	}
}

func analyzeServiceLoad(name string, responseTimes []float64) float64 {
	srv, _ := service.GetServiceByName(name)
	maxRt := srv.Constraints.MaxRespTime
	avgRt := computeAvgResponseTime(responseTimes)
	load := computeLoad(maxRt, avgRt)

	return load
}

func computeAvgResponseTime(responseTimes []float64) float64 {
	sum := 0.0
	avg := 0.0

	for _, rt := range responseTimes {
		sum += rt
	}

	if len(responseTimes) > 0 {
		avg = sum / float64(len(responseTimes))
	}

	return avg
}

func computeLoad(maxRt float64, avgRt float64) float64 {
	// I want the maximum response time
	// to correspond to the 60% of load
	upperBound := maxRt / 0.6
	if avgRt > upperBound {
		avgRt = upperBound
	}

	loadValue := avgRt / upperBound

	log.WithFields(log.Fields{
		"upperBound": upperBound,
		"avgRt":      avgRt,
		"load":       loadValue,
	}).Debugln("Computed load")

	return loadValue
}

func resourcesAvailable(name string) float64 {
	var err error

	nodeMem := node.Config().Resources.TotalMemory
	nodeCpu := node.Config().Resources.TotalCpus
	nodeUsedMem := node.Config().Resources.UsedMemory
	nodeUsedCpu := node.Config().Resources.UsedCpu

	srv, _ := service.GetServiceByName(name)
	srvCpu := getNumberOfCpuFromString(srv.Configuration.CpusetCpus)
	log.WithFields(log.Fields{
		"service": name,
		"cpus":    srvCpu,
	}).Debugln("Service cpu resources")

	var srvMem int64
	if srv.Configuration.Memory != "" {
		srvMem, err = utils.RAMInBytes(srv.Configuration.Memory)
		if err != nil {
			log.WithField("error", err).Warnln("Cannot convert service RAM in Bytes.")
			return 0.0
		}
	} else {
		srvMem = 0
	}

	if nodeCpu < int64(srvCpu) || nodeMem < int64(srvMem) {
		return 0.0
	}

	if (nodeCpu-nodeUsedCpu) < int64(srvCpu) || (nodeMem-nodeUsedMem) < int64(srvMem) {
		return 0.0
	}

	return 1.0
}

func getNumberOfCpuFromString(cpuset string) int64 {
	if cpuset == "" {
		return 0
	}

	return int64(len(strings.Split(cpuset, ",")))
}

func analyzeSystem(analytics *GruAnalytics, stats monitor.GruStats) {
	sysSrvs := []string{}
	for name, _ := range stats.Service {
		sysSrvs = append(sysSrvs, name)
	}

	temp := 0.0
	cpu := stats.System.Cpu
	//TODO compute system mem!!!
	mem := temp
	res := systemResourcesAvailable()
	instances := stats.System.Instances

	health := 1 - ((cpu + mem - res) / 3) //Ok, maybe this is a bit... "mah"...

	sysRes := ResourcesAnalytics{
		cpu,
		mem,
		res,
	}

	SystemAnalytics := SystemAnalytics{
		sysSrvs,
		sysRes,
		instances,
		health,
	}

	gruAnalytics.System = SystemAnalytics
}

func systemResourcesAvailable() float64 {
	totalCpu := float64(node.Config().Resources.TotalCpus)
	totalMemory := float64(node.Config().Resources.TotalMemory)
	usedCpu := float64(node.Config().Resources.UsedCpu)
	usedMemory := float64(node.Config().Resources.UsedMemory)

	cpuRatio := usedCpu / totalCpu
	memRatio := usedMemory / totalMemory
	avgRatio := (cpuRatio + memRatio) / 2

	return (1 - avgRatio)
}

func computeNodeHealth(analytics *GruAnalytics) {
	nServices := len(analytics.Service)
	sumHealth := 0.0
	for _, value := range analytics.Service {
		sumHealth += value.Health
	}
	srvAvgHealth := sumHealth / float64(nServices)

	sysHealth := analytics.System.Health

	totHealth := (srvAvgHealth + sysHealth) / 2

	analytics.Health = totHealth
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
		a, _ := convertDataToAnalytics(data)
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
			if len(ls.Instances.Running) > 0 {
				active = append(active, ls)
			}
		}

		for _, peer := range peers {
			if ps, ok := peer.Service[name]; ok {
				if len(ps.Instances.Running) > 0 {
					active = append(active, ps)
				}
			}
		}

		if len(active) > 1 {
			avgSa = active[0]
			active = active[1:]

			sumLoad := avgSa.Load
			sumCpu := avgSa.Resources.Cpu
			sumMem := avgSa.Resources.Memory
			sumH := avgSa.Health

			for _, actv := range active {

				//LABELS
				sumLoad += actv.Load
				sumCpu += actv.Resources.Cpu
				sumMem += actv.Resources.Memory
				sumH += actv.Health

				//INSTANCES
				avgSa.Instances.All = append(avgSa.Instances.All, actv.Instances.All...)
				avgSa.Instances.Running = append(avgSa.Instances.Running, actv.Instances.Running...)
				avgSa.Instances.Pending = append(avgSa.Instances.Pending, actv.Instances.Pending...)
				avgSa.Instances.Stopped = append(avgSa.Instances.Stopped, actv.Instances.Stopped...)
				avgSa.Instances.Paused = append(avgSa.Instances.Paused, actv.Instances.Paused...)
			}

			total_active := float64(len(active) + 1) // Because I removed the first one from the slice
			avgLoad := sumLoad / total_active
			avgCpu := sumCpu / total_active
			avgMem := sumMem / total_active
			avgH := sumH / total_active

			avgSa.Load = avgLoad
			avgSa.Resources.Cpu = avgCpu
			avgSa.Resources.Memory = avgMem
			avgSa.Health = avgH

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
		sumCpu += peer.System.Resources.Cpu
		sumMem += peer.System.Resources.Memory
		sumH += peer.System.Health
	}

	clstrSrvs = checkAndAppend(clstrSrvs, analytics.System.Services)
	total := float64(len(peers) + 1)
	avgCpu := (analytics.System.Resources.Cpu + sumCpu) / total
	avgMem := (analytics.System.Resources.Memory + sumMem) / total
	avgH := (analytics.System.Health + sumH) / total

	analytics.Cluster.Services = clstrSrvs
	analytics.Cluster.ResourcesAnalytics.Cpu = avgCpu
	analytics.Cluster.ResourcesAnalytics.Memory = avgMem
	analytics.Cluster.Health = avgH
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

func convertAnalyticsToData(analytics GruAnalytics) ([]byte, error) {
	data, err := json.Marshal(analytics)
	if err != nil {
		log.WithField("error", err).Errorln("Error marshaling analytics data")
		return nil, err
	}

	return data, nil
}

func displayAnalyticsOfServices(analytics GruAnalytics) {
	for srv, value := range analytics.Service {
		log.WithFields(log.Fields{
			"service":   srv,
			"cpu":       value.Resources.Cpu,
			"memory":    value.Resources.Memory,
			"resources": value.Resources.Available,
			"load":      value.Load,
			"health":    value.Health,
		}).Infoln("Computed analytics")
	}
}

func GetAnalyzerData() (GruAnalytics, error) {
	analytics := GruAnalytics{}
	dataAnalyics, err := storage.GetData(enum.CLUSTER.ToString(), enum.ANALYTICS)
	if err != nil {
		log.WithField("error", err).Errorln("Cannot retrieve analytics data.")
	} else {
		analytics, err = convertDataToAnalytics(dataAnalyics)
	}

	return analytics, err
}

func convertDataToAnalytics(data []byte) (GruAnalytics, error) {
	analytics := GruAnalytics{}
	err := json.Unmarshal(data, &analytics)
	if err != nil {
		log.WithField("error", err).Errorln("Error unmarshaling analytics data")
	}

	return analytics, err
}
