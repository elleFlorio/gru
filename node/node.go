package node

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/container"
	"github.com/elleFlorio/gru/utils"
)

type Node struct {
	UUID        string      `json:"uuid"`
	Name        string      `json:"name"`
	Constraints Constraints `json:"constraints"`
	Resources   Resources   `json:resources`
}

type Constraints struct {
	CpuMin       float64 `json:"cpumin"`
	CpuMax       float64 `json:"cpumax"`
	MaxInstances int     `json:"maxinstances"` // TODO this will ne removed
}

type Resources struct {
	TotalMemory int64 `json:"totalmemory"`
	TotalCpus   int64 `json:"totalcpus"`
	usedMemory  int64 `json:"usedmemory"`
	usedCpu     int64 `json:"usedcpu"`
}

var node Node

func LoadNodeConfig(filename string) error {
	node.UUID, _ = utils.GenerateUUID()

	log.WithField("status", "start").Infoln("Node configuration loading")
	defer log.WithFields(log.Fields{
		"status": "done",
		"UUID":   node.UUID,
	}).Infoln("Node configuration loading")

	tmp, err := ioutil.ReadFile(filename)
	if err != nil {
		log.WithField("error", err).Errorln("Error reading node configuration file")
		return err
	}
	err = json.Unmarshal(tmp, &node)
	if err != nil {
		log.WithField("error", err).Errorln("Error unmarshaling node configuration file")
		return err
	}

	return nil
}

//Use overcommit ratio?
func ComputeTotalResources() {
	info, err := container.Docker().Client.Info()
	if err != nil {
		log.WithField("error", err).Errorln("Error reading total resources")
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
		cData, err := container.Docker().Client.InspectContainer(c.Id)
		if err != nil {
			return 0, err
		}

		cpus += cData.Config.CpuShares
	}

	node.Resources.TotalCpus = cpus

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

	node.Resources.TotalMemory = memory

	return memory, nil
}

func Config() Node {
	return node
}

func UpdateNode(newNode Node) {
	node = newNode
}
