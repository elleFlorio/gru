package autonomic

import (
//"github.com/samalba/dockerclient"
)

type dckrData map[string]containerData

type analyzeData struct {
	data dckrData
}

type containerData struct {
	Info struct {
		Id    string
		Image string
	}
	Cpu struct {
		ContPerc float64
		ContJif  uint64
		SysJif   uint64
	}
}

func (analyze analyzeData) getContainerIds() []string {
	ids := make([]string, len(analyze.data), len(analyze.data))
	idx := 0
	for key := range analyze.data {
		ids[idx] = key
		idx++
	}
	return ids
}

func (analyze analyzeData) getAllData() []containerData {
	dt := make([]containerData, len(analyze.data), len(analyze.data))
	idx := 0
	for _, value := range analyze.data {
		dt[idx] = value
		idx++
	}
	return dt
}

func (analyze analyzeData) groupByImage() map[string][]containerData {
	group := make(map[string][]containerData)
	for _, cd := range analyze.getAllData() {
		group[cd.Info.Image] = append(group[cd.Info.Image], cd)
	}
	return group
}

func (analyze analyzeData) isEmpty() bool {
	return len(analyze.data) == 0
}
