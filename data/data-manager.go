package data

import (
	"encoding/json"
	"errors"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/service"
	"github.com/elleFlorio/gru/storage"
)

func SaveStats(stats GruStats) {
	err := saveData(stats, enum.STATS)
	if err != nil {
		log.WithField("err", err).Debugln("Cannot convert stats to data")
	}
}

func SaveAnalytics(analytics GruAnalytics) {
	err := saveData(analytics, enum.ANALYTICS)
	if err != nil {
		log.WithField("err", err).Debugln("Cannot convert stats to data")
	}
}

func SavePolicy(policy Policy) {
	err := saveData(policy, enum.POLICIES)
	if err != nil {
		log.WithField("err", err).Debugln("Cannot convert stats to data")
	}
}

func SaveInfo(info GruInfo) {
	err := saveData(info, enum.INFO)
	if err != nil {
		log.WithField("err", err).Debugln("Cannot convert stats to data")
	}
}

func saveData(data interface{}, dataType enum.Datatype) error {
	var encoded []byte
	var err error
	switch dataType {
	case enum.STATS:
		stats := data.(GruStats)
		encoded, err = json.Marshal(stats)
		if err != nil {
			return err
		}
	case enum.ANALYTICS:
		analytics := data.(GruAnalytics)
		encoded, err = json.Marshal(analytics)
		if err != nil {
			return err
		}
	case enum.POLICIES:
		policy := data.(Policy)
		encoded, err = json.Marshal(policy)
		if err != nil {
			return err
		}
	case enum.INFO:
		info := data.(GruInfo)
		encoded, err = json.Marshal(info)
		if err != nil {
			return err
		}
	default:
		return errors.New("Cannot save data: unknown data type")
	}

	storage.StoreClusterData(encoded, dataType)

	return nil
}

func GetStats() (GruStats, error) {
	stats, err := getData(enum.STATS)
	if err != nil {
		log.WithField("err", err).Warnln("Cannot get stats data")
		return GruStats{}, err
	}

	return stats.(GruStats), nil
}

func GetAnalytics() (GruAnalytics, error) {
	analytics, err := getData(enum.ANALYTICS)
	if err != nil {
		log.WithField("err", err).Warnln("Cannot get analytics data")
		return GruAnalytics{}, err
	}

	return analytics.(GruAnalytics), nil
}

func GetPolicy() (Policy, error) {
	policy, err := getData(enum.POLICIES)
	if err != nil {
		log.WithField("err", err).Warnln("Cannot get policy data")
		return Policy{}, err
	}

	return policy.(Policy), nil
}

func GetInfo() (GruInfo, error) {
	info, err := getData(enum.INFO)
	if err != nil {
		log.WithField("err", err).Warnln("Cannot get info data")
		return GruInfo{}, err
	}

	return info.(GruInfo), nil
}

func getData(dataType enum.Datatype) (interface{}, error) {
	var data interface{}
	switch dataType {
	case enum.STATS:
		data = GruStats{}
		dataStats, err := storage.GetClusterData(dataType)
		if err != nil {
			return nil, err
		} else {
			data, err = ByteToStats(dataStats)
			if err != nil {
				return nil, err
			}
		}
	case enum.ANALYTICS:
		data = GruAnalytics{}
		dataAnalytics, err := storage.GetClusterData(dataType)
		if err != nil {
			return nil, err
		} else {
			data, err = ByteToAnalytics(dataAnalytics)
			if err != nil {
				return nil, err
			}
		}
	case enum.POLICIES:
		data = Policy{}
		dataPolicy, err := storage.GetClusterData(dataType)
		if err != nil {
			return nil, err
		} else {
			data, err = ByteToPolicy(dataPolicy)
			if err != nil {
				return nil, err
			}
		}
	case enum.INFO:
		data = GruInfo{}
		dataInfo, err := storage.GetClusterData(dataType)
		if err != nil {
			return nil, err
		} else {
			data, err = ByteToInfo(dataInfo)
			if err != nil {
				return nil, err
			}
		}
	}

	return data, nil

}

func ByteToStats(data []byte) (GruStats, error) {
	stats := GruStats{}
	err := json.Unmarshal(data, &stats)
	if err != nil {
		log.WithField("err", err).Warnln("Cannot convert byte to stats")
		return GruStats{}, err
	}

	return stats, nil

}

func ByteToAnalytics(data []byte) (GruAnalytics, error) {
	analytics := GruAnalytics{}
	err := json.Unmarshal(data, &analytics)
	if err != nil {
		log.WithField("err", err).Warnln("Cannot conver byte to analytics")
		return GruAnalytics{}, err
	}

	return analytics, nil

}

func ByteToPolicy(data []byte) (Policy, error) {
	policy := Policy{}
	err := json.Unmarshal(data, &policy)
	if err != nil {
		log.WithField("err", err).Warnln("Cannot conver byte to policy")
		return Policy{}, err
	}

	return policy, nil

}

func ByteToInfo(data []byte) (GruInfo, error) {
	info := GruInfo{}
	err := json.Unmarshal(data, &info)
	if err != nil {
		log.WithField("err", err).Warnln("Cannot conver byte to info")
		return GruInfo{}, err
	}

	return info, nil

}

func MergeInfo(toMerge []GruInfo) GruInfo {
	loadAvg := 0.0
	cpuAvg := 0.0
	memAvg := 0.0
	resourcesAvg := 0.0

	merged := GruInfo{
		Service: make(map[string]ServiceInfo),
	}

	for _, name := range service.List() {
		counter := 0.0
		for _, info := range toMerge {
			if srv, ok := info.Service[name]; ok {
				loadAvg += srv.Load
				cpuAvg += srv.Cpu
				memAvg += srv.Memory
				resourcesAvg += srv.Resources

				counter++
			}
		}

		if counter > 0 {
			loadAvg /= counter
			cpuAvg /= counter
			memAvg /= counter
			resourcesAvg /= counter

			mergedService := ServiceInfo{
				Load:      loadAvg,
				Cpu:       cpuAvg,
				Memory:    memAvg,
				Resources: resourcesAvg,
				Active:    true,
			}

			merged.Service[name] = mergedService
		}
	}

	cpuAvg = 0.0
	memAvg = 0.0
	resourcesAvg = 0.0
	healthAvg := 0.0
	activeServices := []string{}
	for _, info := range toMerge {
		cpuAvg += info.System.Cpu
		memAvg += info.System.Memory
		healthAvg += info.System.Health
		activeServices = checkAndAppend(activeServices, info.System.ActiveServices)
	}

	lenght := float64(len(toMerge))
	cpuAvg /= lenght
	memAvg /= lenght
	healthAvg /= lenght

	mergedSystem := SystemInfo{
		Cpu:            cpuAvg,
		Memory:         memAvg,
		Health:         healthAvg,
		ActiveServices: activeServices,
	}

	merged.System = mergedSystem

	return merged
}

func checkAndAppend(list []string, toAppend []string) []string {
	if len(list) == 0 {
		return append(list, toAppend...)
	}

	var isPresent bool
	for _, elem := range toAppend {
		isPresent = false
		for _, present := range list {
			if elem == present {
				isPresent = true
			}
		}

		if !isPresent {
			list = append(list, elem)
		}
	}

	return list
}
