package autonomic

import (
	//log "github.com/Sirupsen/logrus"
	"github.com/samalba/dockerclient"
)

type dckrStats map[string]dockerclient.Stats
type dckrInfo map[string]dockerclient.ContainerInfo

type monitorData struct {
	stats dckrStats
	info  dckrInfo
}

func (monitor monitorData) getContainerIds() []string {
	ids := make([]string, len(monitor.info), len(monitor.info))
	idx := 0
	for key := range monitor.info {
		ids[idx] = key
		idx++
	}
	return ids
}

func (monitor monitorData) getAllStats() []dockerclient.Stats {
	st := make([]dockerclient.Stats, len(monitor.stats), len(monitor.stats))
	idx := 0
	for _, value := range monitor.stats {
		st[idx] = value
		idx++
	}
	return st
}

func (monitor monitorData) getAllInfo() []dockerclient.ContainerInfo {
	info := make([]dockerclient.ContainerInfo, len(monitor.info), len(monitor.info))
	idx := 0
	for _, value := range monitor.info {
		info[idx] = value
		idx++
	}
	return info
}

func (monitor monitorData) isEmpty() bool {
	return len(monitor.stats) == 0 && len(monitor.info) == 0
}

func (monitor *monitorData) removeData(id string) {
	delete(monitor.stats, id)
	delete(monitor.info, id)
}
