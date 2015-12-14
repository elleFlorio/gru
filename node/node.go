package node

import (
	"strings"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/container"
	"github.com/elleFlorio/gru/network"
	"github.com/elleFlorio/gru/service"
	"github.com/elleFlorio/gru/utils"
)

var node Node

func CreateNode(name string) {
	node_UUID, err := utils.GenerateUUID()
	if err != nil {
		log.WithField("err", err).Errorln("Error generating node UUID")
	}
	node_address := "http://" + network.Config().IpAddress + ":" + network.Config().Port

	config := Config{node_UUID, name, node_address, ""}
	node = Node{
		Configuration: config,
		Active:        false,
	}

	computeTotalResources()
}

func computeTotalResources() {
	info, err := container.Docker().Client.Info()
	if err != nil {
		log.WithField("err", err).Errorln("Error reading total resources")
		return
	}
	node.Resources.TotalCpus = info.NCPU
	node.Resources.TotalMemory = info.MemTotal
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

	node.Resources.UsedCpu = cpus

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

	node.Resources.UsedMemory = memory

	return memory, nil
}

func GetNode() Node {
	return node
}

func ToggleActiveNode() {
	node.Active = !node.Active
}

func UpdateNode(newNode Node) {
	node = newNode
}
