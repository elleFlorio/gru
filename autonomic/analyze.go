package autonomic

import (
	log "github.com/Sirupsen/logrus"
	//"github.com/samalba/dockerclient"
)

var a_data analyzeData = analyzeData{
	make(dckrData),
}

func analyze(m_data *monitorData) *analyzeData {
	//Analyze stuff
	log.Debugln("I'm analyzing")

	for _, id := range m_data.getContainerIds() {
		stats := a_data.data[id]
		stats.Info.Id = id
		stats.Info.Image = m_data.info[id].Image
		computeCpuUsagePerc(id, m_data, stats)
	}

	return &a_data
}

/*
ref: http://stackoverflow.com/questions/1420426/calculating-cpu-usage-of-a-process-in-linux
*/
func computeCpuUsagePerc(id string, m_data *monitorData, stats containerData) {
	tot_old := stats.Cpu.ContJif
	sys_old := stats.Cpu.SysJif
	tot := m_data.stats[id].CpuStats.CpuUsage.TotalUsage
	sys := m_data.stats[id].CpuStats.SystemUsage

	if !(tot_old == 0) || !(sys_old == 0) {
		totCpuPerc := 100 * float64(tot-tot_old) / float64(sys-sys_old)
		stats.Cpu.ContPerc = totCpuPerc
	}
	stats.Cpu.ContJif = tot
	stats.Cpu.SysJif = sys
	a_data.data[id] = stats

	log.WithFields(log.Fields{
		"Name": m_data.info[id].Name,
		"CPU":  a_data.data[id].Cpu.ContPerc,
	}).Debugln("Computed CPU usage")
}
