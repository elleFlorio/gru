package analyzer

import (
	"encoding/json"
	"errors"
	"fmt"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/data"
	"github.com/elleFlorio/gru/enum"
	res "github.com/elleFlorio/gru/resources"
	"github.com/elleFlorio/gru/service"
	"github.com/elleFlorio/gru/storage"
)

const overcommitratio float64 = 0.0

var (
	gruAnalytics          data.GruAnalytics
	ErrNoRunningInstances error = errors.New("No active instance to analyze")
)

func init() {
	gruAnalytics = data.GruAnalytics{
		Service: make(map[string]data.ServiceAnalytics),
	}
}

func Run(stats data.GruStats) data.GruAnalytics {
	log.WithField("status", "init").Debugln("Gru Analyzer")
	defer log.WithField("status", "done").Debugln("Gru Analyzer")

	if len(stats.Service) == 0 {
		log.WithField("err", "No services stats").Warnln("Cannot compute analytics.")
	} else {
		updateNodeResources()
		analyzeServices(&gruAnalytics, stats)
		analyzeSystem(&gruAnalytics, stats)
		computeNodeHealth(&gruAnalytics)
		analyzeCluster(&gruAnalytics)
		data.SaveAnalytics(gruAnalytics)
		displayAnalyticsOfServices(gruAnalytics)
	}

	return gruAnalytics
}

func updateNodeResources() {
	res.ComputeUsedResources()

	log.WithFields(log.Fields{
		"totalcpu": res.GetResources().CPU.Total,
		"usedcpu":  res.GetResources().CPU.Used,
		"totalmem": res.GetResources().Memory.Total,
		"usedmem":  res.GetResources().Memory.Used,
	}).Debugln("Updated node resources")
}

func analyzeServices(analytics *data.GruAnalytics, stats data.GruStats) {
	for name, value := range stats.Service {
		load := analyzeServiceLoad(name, value.Metrics.ResponseTime)
		cpu := value.Cpu.Tot
		mem := value.Memory.Tot
		resAvailable := res.AvailableResourcesService(name)

		srv, _ := service.GetServiceByName(name)
		instances := srv.Instances

		health := 1 - ((load + mem + cpu - resAvailable) / 4) //I don't like this...

		srvRes := data.ResourcesAnalytics{
			cpu,
			mem,
			resAvailable,
		}

		srvAnalytics := data.ServiceAnalytics{
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
	// to correspond to the 80% of load
	upperBound := maxRt / 0.8
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

func analyzeSystem(analytics *data.GruAnalytics, stats data.GruStats) {
	sysSrvs := []string{}
	for name, _ := range stats.Service {
		sysSrvs = append(sysSrvs, name)
	}

	temp := 0.0
	cpu := stats.System.Cpu
	//TODO compute system mem!!!
	mem := temp
	resources := res.AvailableResources()
	instances := *cfg.GetNodeInstances()

	health := 1 - ((cpu + mem - resources) / 3) //Ok, maybe this is a bit... "mah"...

	sysRes := data.ResourcesAnalytics{
		cpu,
		mem,
		resources,
	}

	systemAnalytics := data.SystemAnalytics{
		sysSrvs,
		sysRes,
		instances,
		health,
	}

	gruAnalytics.System = systemAnalytics
}

func computeNodeHealth(analytics *data.GruAnalytics) {
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

func analyzeCluster(analytics *data.GruAnalytics) {
	peers := getPeersAnalytics()
	computeServicesAvg(peers, analytics)
	computeClusterAvg(peers, analytics)
}

func getPeersAnalytics() []data.GruAnalytics {
	peers := make([]data.GruAnalytics, 0)
	dataAn, _ := storage.GetAllData(enum.ANALYTICS)
	delete(dataAn, enum.CLUSTER.ToString())
	for _, data := range dataAn {
		a, _ := convertDataToAnalytics(data)
		peers = append(peers, a)
	}

	return peers
}

func computeServicesAvg(peers []data.GruAnalytics, analytics *data.GruAnalytics) {
	avg := make(map[string]data.ServiceAnalytics)

	for _, name := range service.List() {
		active := []data.ServiceAnalytics{}
		var avgSa data.ServiceAnalytics
		isLocal := false

		if ls, ok := analytics.Service[name]; ok {
			if len(ls.Instances.Running) > 0 {
				active = append(active, ls)
				isLocal = true
				log.WithField("service", name).Debugln("Local service running")
			}
		}

		for _, peer := range peers {
			if ps, ok := peer.Service[name]; ok {
				if len(ps.Instances.Running) > 0 {
					active = append(active, ps)
				}
			}
		}

		log.WithFields(log.Fields{
			"service": name,
			"total":   len(active),
		}).Debugln("Active services")

		if len(active) > 0 {
			if !isLocal {
				// I don't want to merge instance status. This may
				// cause some problems of synchronization with
				// other peers
				active[0].Instances = cfg.ServiceStatus{}
				// I don't have to take the resources available for that service
				// in the peer node, I need to compute the local ones
				active[0].Resources.Available = res.AvailableResourcesService(name)
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

				sumLoad += actv.Load
				sumCpu += actv.Resources.Cpu
				sumMem += actv.Resources.Memory
				sumH += actv.Health
			}

			total_active := float64(len(active) + 1) // Because I removed the first one from the slice
			avgLoad := sumLoad / total_active
			avgCpu := sumCpu / total_active
			avgMem := sumMem / total_active
			avgH := sumH / total_active

			log.WithFields(log.Fields{
				"service":      name,
				"sumLoad":      sumLoad,
				"sumCpu":       sumCpu,
				"total_active": total_active,
			}).Debugln("Active services")

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

func computeClusterAvg(peers []data.GruAnalytics, analytics *data.GruAnalytics) {
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
	analytics.Cluster.Resources.Cpu = avgCpu
	analytics.Cluster.Resources.Memory = avgMem
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
// func saveAnalytics(analytics data.GruAnalytics) error {
// 	data, err := convertAnalyticsToData(analytics)
// 	if err != nil {
// 		log.WithField("err", err).Errorln("Cannot convert analytics to data")
// 		return err
// 	} else {
// 		storage.StoreData(enum.CLUSTER.ToString(), data, enum.ANALYTICS)
// 	}

// 	return nil
// }

// func convertAnalyticsToData(analytics data.GruAnalytics) ([]byte, error) {
// 	data, err := json.Marshal(analytics)
// 	if err != nil {
// 		log.WithField("err", err).Errorln("Error marshaling analytics data")
// 		return nil, err
// 	}

// 	return data, nil
// }

func displayAnalyticsOfServices(analytics data.GruAnalytics) {
	for srv, value := range analytics.Service {
		log.WithFields(log.Fields{
			"service":   srv,
			"cpu":       fmt.Sprintf("%.2f", value.Resources.Cpu),
			"memory":    fmt.Sprintf("%.2f", value.Resources.Memory),
			"resources": fmt.Sprintf("%.2f", value.Resources.Available),
			"load":      fmt.Sprintf("%.2f", value.Load),
			"health":    fmt.Sprintf("%.2f", value.Health),
		}).Infoln("Analytics computed: ", srv)
	}
}

// func GetAnalyzerData() (data.GruAnalytics, error) {
// 	analytics := data.GruAnalytics{}
// 	dataAnalyics, err := storage.GetData(enum.CLUSTER.ToString(), enum.ANALYTICS)
// 	if err != nil {
// 		log.WithField("err", err).Warnln("Cannot retrieve analytics data")
// 	} else {
// 		analytics, err = convertDataToAnalytics(dataAnalyics)
// 	}

// 	return analytics, err
// }

//TODO to be removed after data information integration
func convertDataToAnalytics(dataByte []byte) (data.GruAnalytics, error) {
	analytics := data.GruAnalytics{}
	err := json.Unmarshal(dataByte, &analytics)
	if err != nil {
		log.WithField("err", err).Warnln("Error converting data to analytics")
		return data.GruAnalytics{}, err
	}

	return analytics, nil
}

func analyzeSharedData(analytics *data.GruAnalytics) data.Shared {
	myShared := computeLocalShared(analytics)
	data.SaveSharedLocal(myShared)
	clusterData := computeClusterData(myShared)
	data.SaveSharedCluster(clusterData)
	return clusterData

}

func computeLocalShared(analytics *data.GruAnalytics) data.Shared {
	myShared := data.Shared{Service: make(map[string]data.ServiceShared)}

	for srv, value := range analytics.Service {
		mySrvShared := data.ServiceShared{
			Load:      value.Load,
			Cpu:       value.Resources.Cpu,
			Memory:    value.Resources.Memory,
			Resources: value.Resources.Available,
			Active:    true,
		}

		myShared.Service[srv] = mySrvShared
	}

	myShared.System.Cpu = analytics.System.Resources.Cpu
	myShared.System.Memory = analytics.System.Resources.Memory
	myShared.System.Health = analytics.System.Health
	myShared.System.ActiveServices = analytics.System.Services

	return myShared
}

func computeClusterData(myShared data.Shared) data.Shared {
	sharedData, err := data.GetSharedCluster()
	if err != nil {
		log.Debugln("Cannot compute cluster data")
		return myShared
	}

	toMerge := []data.Shared{myShared, sharedData}
	clusterData, err := data.MergeShared(toMerge)
	if err != nil {
		return myShared
	}

	return clusterData
}
