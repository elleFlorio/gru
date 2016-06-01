package analyzer

import (
	"errors"
	"fmt"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/data"
	res "github.com/elleFlorio/gru/resources"
	"github.com/elleFlorio/gru/service"
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

func Run(stats data.GruStats) data.Shared {
	log.WithField("status", "init").Debugln("Gru Analyzer")
	defer log.WithField("status", "done").Debugln("Gru Analyzer")

	sharedData := data.Shared{Service: make(map[string]data.ServiceShared)}

	if len(stats.Service) == 0 {
		log.WithField("err", "No services stats").Warnln("Cannot compute analytics.")
	} else {
		updateNodeResources()
		analyzeServices(&gruAnalytics, stats)
		analyzeSystem(&gruAnalytics, stats)
		computeNodeHealth(&gruAnalytics)
		data.SaveAnalytics(gruAnalytics)
		sharedData = analyzeSharedData(&gruAnalytics)

		displayAnalyticsOfServices(gruAnalytics)
		displaySharedInformation(sharedData)
	}

	return sharedData
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
		cpu := value.Cpu.Avg
		mem := value.Memory.Avg
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
	lowerbound := maxRt / 2.0
	if avgRt < lowerbound {
		return 0.0
	}

	upperBound := maxRt
	if avgRt > upperBound {
		return 1.0
	}

	loadValue := (avgRt - lowerbound) / (upperBound - lowerbound)

	log.WithField("load", loadValue).Debugln("Computed load")

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
			Active:    isServiceActive(value.Instances),
		}

		myShared.Service[srv] = mySrvShared
	}

	myShared.System.Cpu = analytics.System.Resources.Cpu
	myShared.System.Memory = analytics.System.Resources.Memory
	myShared.System.Health = analytics.System.Health
	myShared.System.ActiveServices = analytics.System.Services

	return myShared
}

func isServiceActive(status cfg.ServiceStatus) bool {
	return (len(status.Pending) + len(status.Running)) > 0
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

func displayAnalyticsOfServices(analytics data.GruAnalytics) {
	for srv, value := range analytics.Service {
		log.WithFields(log.Fields{
			"service":   srv,
			"cpu":       fmt.Sprintf("%.2f", value.Resources.Cpu),
			"memory":    fmt.Sprintf("%.2f", value.Resources.Memory),
			"resources": fmt.Sprintf("%.2f", value.Resources.Available),
			"load":      fmt.Sprintf("%.2f", value.Load),
			"health":    fmt.Sprintf("%.2f", value.Health),
		}).Infoln("Analytics computed")
	}
}

func displaySharedInformation(clusterData data.Shared) {
	for srv, value := range clusterData.Service {
		log.WithFields(log.Fields{
			"service":   srv,
			"cpu":       fmt.Sprintf("%.2f", value.Cpu),
			"memory":    fmt.Sprintf("%.2f", value.Memory),
			"resources": fmt.Sprintf("%.2f", value.Resources),
			"load":      fmt.Sprintf("%.2f", value.Load),
		}).Infoln("Cluster shared data")
	}
}
