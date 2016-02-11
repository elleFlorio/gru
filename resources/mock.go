package resources

import (
	"github.com/elleFlorio/gru/utils"
)

func CreateMockResources(totCpu int64, totMem string, usedCpu int64, usedMem string) {
	setResources(totCpu, totMem, usedCpu, usedMem)
}

func setResources(totCpu int64, totMem string, usedCpu int64, usedMem string) {
	totMemB, _ := utils.RAMInBytes(totMem)
	usedMemB, _ := utils.RAMInBytes(usedMem)
	resources.Memory.Total = totMemB
	resources.Memory.Used = usedMemB
	resources.CPU.Total = totCpu
	resources.CPU.Used = usedCpu
}
