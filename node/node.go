package node

import (
	"strings"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/container"
	"github.com/elleFlorio/gru/network"
	"github.com/elleFlorio/gru/service"
	"github.com/elleFlorio/gru/utils"
)

func CreateNode(name string) {
	node_UUID, err := utils.GenerateUUID()
	if err != nil {
		log.WithField("err", err).Errorln("Error generating node UUID")
	}
	node_address := "http://" + network.Config().IpAddress + ":" + network.Config().Port
	config := cfg.NodeConfig{node_UUID, name, node_address, "", ""}
	node := cfg.Node{
		Configuration: config,
		Active:        false,
	}
	cfg.SetNode(node)

	computeTotalResources()
}

func computeTotalResources() {
	info, err := container.Docker().Client.Info()
	if err != nil {
		log.WithField("err", err).Errorln("Error reading total resources")
		return
	}
	cfg.GetNodeResources().TotalCpus = info.NCPU
	cfg.GetNodeResources().TotalMemory = info.MemTotal
}

func UsedCpus() (int64, error) {
	var cpus int64

	containers, err := container.Docker().Client.ListContainers(false, false, "")
	if err != nil {
		return 0, err
	}

	for _, c := range containers {
		if _, err := service.GetServiceByImage(c.Image); err == nil {
			cData, err := container.Docker().Client.InspectContainer(c.Id)
			if err != nil {
				return 0, err
			}
			cpuset := strings.Split(cData.HostConfig.CpusetCpus, ",")
			cpus += int64(len(cpuset))
		}
	}

	cfg.GetNodeResources().UsedCpu = cpus

	return cpus, nil
}

func UsedMemory() (int64, error) {
	var memory int64

	containers, err := container.Docker().Client.ListContainers(false, false, "")
	if err != nil {
		return 0, err
	}

	for _, c := range containers {
		cData, err := container.Docker().Client.InspectContainer(c.Id)
		if err != nil {
			return 0, err
		}

		memory += cData.Config.Memory
	}

	cfg.GetNodeResources().UsedMemory = memory

	return memory, nil
}
