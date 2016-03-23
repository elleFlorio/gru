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
	err := saveData(stats, enum.STATS, enum.LOCAL)
	if err != nil {
		log.WithField("err", err).Debugln("Cannot convert stats to data")
	}
}

func SaveAnalytics(analytics GruAnalytics) {
	err := saveData(analytics, enum.ANALYTICS, enum.LOCAL)
	if err != nil {
		log.WithField("err", err).Debugln("Cannot convert stats to data")
	}
}

func SavePolicy(policy Policy) {
	err := saveData(policy, enum.POLICIES, enum.LOCAL)
	if err != nil {
		log.WithField("err", err).Debugln("Cannot convert stats to data")
	}
}

func SaveShared(info Shared) {
	err := saveData(info, enum.SHARED, enum.CLUSTER)
	if err != nil {
		log.WithField("err", err).Debugln("Cannot convert stats to data")
	}
}

func saveData(data interface{}, dataType enum.Datatype, dataOwner enum.DataOwner) error {
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
	case enum.SHARED:
		info := data.(Shared)
		encoded, err = json.Marshal(info)
		if err != nil {
			return err
		}
	default:
		return errors.New("Cannot save data: unknown data type")
	}

	storage.StoreData(dataOwner.ToString(), encoded, dataType)

	return nil
}

func GetStats() (GruStats, error) {
	stats, err := getData(enum.STATS, enum.LOCAL)
	if err != nil {
		log.WithField("err", err).Warnln("Cannot get stats data")
		return GruStats{}, err
	}

	return stats.(GruStats), nil
}

func GetAnalytics() (GruAnalytics, error) {
	analytics, err := getData(enum.ANALYTICS, enum.LOCAL)
	if err != nil {
		log.WithField("err", err).Warnln("Cannot get analytics data")
		return GruAnalytics{}, err
	}

	return analytics.(GruAnalytics), nil
}

func GetPolicy() (Policy, error) {
	policy, err := getData(enum.POLICIES, enum.LOCAL)
	if err != nil {
		log.WithField("err", err).Warnln("Cannot get policy data")
		return Policy{}, err
	}

	return policy.(Policy), nil
}

func GetShared() (Shared, error) {
	info, err := getData(enum.SHARED, enum.CLUSTER)
	if err != nil {
		log.WithField("err", err).Warnln("Cannot get info data")
		return Shared{}, err
	}

	return info.(Shared), nil
}

func getData(dataType enum.Datatype, dataOwner enum.DataOwner) (interface{}, error) {
	var data interface{}
	switch dataType {
	case enum.STATS:
		data = GruStats{}
		dataStats, err := storage.GetData(dataOwner.ToString(), dataType)
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
		dataAnalytics, err := storage.GetData(dataOwner.ToString(), dataType)
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
		dataPolicy, err := storage.GetData(dataOwner.ToString(), dataType)
		if err != nil {
			return nil, err
		} else {
			data, err = ByteToPolicy(dataPolicy)
			if err != nil {
				return nil, err
			}
		}
	case enum.SHARED:
		data = Shared{}
		dataInfo, err := storage.GetData(dataOwner.ToString(), dataType)
		if err != nil {
			return nil, err
		} else {
			data, err = ByteToShared(dataInfo)
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

func ByteToShared(data []byte) (Shared, error) {
	info := Shared{}
	err := json.Unmarshal(data, &info)
	if err != nil {
		log.WithField("err", err).Warnln("Cannot conver byte to info")
		return Shared{}, err
	}

	return info, nil

}

func MergeShared(toMerge []Shared) (Shared, error) {
	if len(toMerge) < 1 {
		return Shared{}, errors.New("No shared data to merge")
	}

	if len(toMerge) == 1 {
		return toMerge[0], nil
	}

	loadAvg := 0.0
	cpuAvg := 0.0
	memAvg := 0.0
	resourcesAvg := 0.0

	merged := Shared{
		Service: make(map[string]ServiceShared),
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

			mergedService := ServiceShared{
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

	mergedSystem := SystemShared{
		Cpu:            cpuAvg,
		Memory:         memAvg,
		Health:         healthAvg,
		ActiveServices: activeServices,
	}

	merged.System = mergedSystem

	return merged, nil
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
