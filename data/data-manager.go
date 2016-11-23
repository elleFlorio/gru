package data

import (
	"encoding/json"
	"errors"
	"sync"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/enum"
	srv "github.com/elleFlorio/gru/service"
	"github.com/elleFlorio/gru/storage"
	"github.com/elleFlorio/gru/utils"
)

var (
	friendsData utils.LkList
	m_friends   = &sync.RWMutex{}
)

func InitializeFriendsData(limit int) {
	friendsData = utils.CreateLkList(limit)
}

func AddFriendData(friend string, value Shared) {
	defer m_friends.Unlock()

	m_friends.Lock()
	friendsData.PushValue(friend, value)
}

func GetFriendsData() []Shared {
	defer m_friends.RUnlock()

	values := make([]Shared, friendsData.Limit)
	m_friends.RLock()
	for _, value := range friendsData.GetValues() {
		values = append(values, value.(Shared))
	}
	return values
}

func SaveStats(stats GruStats) {
	err := saveData(stats, enum.STATS, enum.LOCAL)
	if err != nil {
		log.WithField("err", err).Debugln("Cannot convert stats to data")
	}
}

func SaveAnalytics(analytics GruAnalytics) {
	err := saveData(analytics, enum.ANALYTICS, enum.LOCAL)
	if err != nil {
		log.WithField("err", err).Debugln("Cannot convert analytics to data")
	}
}

func SavePolicy(policy Policy) {
	err := saveData(policy, enum.POLICIES, enum.LOCAL)
	if err != nil {
		log.WithField("err", err).Debugln("Cannot convert policy to data")
	}
}

func SaveSharedLocal(info Shared) {
	err := saveData(info, enum.SHARED, enum.LOCAL)
	if err != nil {
		log.WithField("err", err).Debugln("Cannot convert local shared to data")
	}
}

func SaveSharedCluster(info Shared) {
	err := saveData(info, enum.SHARED, enum.CLUSTER)
	if err != nil {
		log.WithField("err", err).Debugln("Cannot convert cluster shared to data")
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

func GetSharedLocal() (Shared, error) {
	info, err := getData(enum.SHARED, enum.LOCAL)
	if err != nil {
		log.WithField("err", err).Warnln("Cannot get shared local data")
		return Shared{}, err
	}

	return info.(Shared), nil
}

func GetSharedCluster() (Shared, error) {
	info, err := getData(enum.SHARED, enum.CLUSTER)
	if err != nil {
		log.WithField("err", err).Warnln("Cannot get shared cluster data")
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

func SharedToByte(data Shared) []byte {
	encoded, err := json.Marshal(data)
	if err != nil {
		log.WithField("err", err).Warnln("Cannot convert shared data to byte")
		return []byte{}
	}

	return encoded
}

func MergeShared(toMerge []Shared) (Shared, error) {
	if len(toMerge) < 1 {
		return Shared{}, errors.New("No shared data to merge")
	}

	if len(toMerge) == 1 {
		return toMerge[0], nil
	}

	merged := Shared{
		Service: make(map[string]ServiceShared),
	}

	for _, name := range srv.List() {
		srvMerged := ServiceShared{}
		baseValues := make(map[string][]float64)
		userValues := make(map[string][]float64)
		for _, data := range toMerge {
			if data.Service[name].Active {
				srvMerged.Active = true

				for analytics, value := range data.Service[name].Data.BaseShared {
					baseValues[analytics] = append(baseValues[analytics], value)
				}

				for analytics, value := range data.Service[name].Data.UserShared {
					userValues[analytics] = append(userValues[analytics], value)
				}
			}
		}

		baseMerged := make(map[string]float64, len(baseValues))
		for analytics, values := range baseValues {
			baseMerged[analytics] = utils.Mean(values)
		}

		userMerged := make(map[string]float64, len(userValues))
		for analytics, values := range userValues {
			userMerged[analytics] = utils.Mean(values)
		}

		srvMerged.Data.BaseShared = baseMerged
		srvMerged.Data.UserShared = userMerged
		merged.Service[name] = srvMerged
	}

	sysMerged := SystemShared{}
	baseValues := make(map[string][]float64)
	userValues := make(map[string][]float64)
	for _, data := range toMerge {

		for analytics, value := range data.System.Data.BaseShared {
			baseValues[analytics] = append(baseValues[analytics], value)
		}

		for analytics, value := range data.System.Data.UserShared {
			userValues[analytics] = append(userValues[analytics], value)
		}

		sysMerged.ActiveServices = checkAndAppend(sysMerged.ActiveServices, data.System.ActiveServices)
	}

	baseMerged := make(map[string]float64, len(baseValues))
	for analytics, values := range baseValues {
		baseMerged[analytics] = utils.Mean(values)
	}

	userMerged := make(map[string]float64, len(userValues))
	for analytics, values := range userValues {
		userMerged[analytics] = utils.Mean(values)
	}

	sysMerged.Data.BaseShared = baseMerged
	sysMerged.Data.UserShared = userMerged
	merged.System = sysMerged

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
