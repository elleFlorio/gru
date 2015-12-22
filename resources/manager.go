package resources

import (
	"errors"
	"runtime"
	"strconv"
	"strings"
	"sync"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/container"
	"github.com/elleFlorio/gru/service"
)

var (
	resources Resource
	mutex_cpu = sync.RWMutex{}

	ErrTooManyCores    = errors.New("Number of requested cores exceed the total number of available ones")
	ErrWrongCoreNumber = errors.New("The number of the specified core is not correct")
)

func init() {
	resources = Resource{}
	resources.CPU.Cores = make(map[int]bool)
}

func Initialize() {
	computeTotalResources()
}

func computeTotalResources() {
	info, err := container.Docker().Client.Info()
	if err != nil {
		log.WithField("err", err).Errorln("Error reading total resources")
		return
	}

	resources.CPU.Total = info.NCPU
	resources.Memory.Total = info.MemTotal

	for i := 0; i < int(resources.CPU.Total); i++ {
		resources.CPU.Cores[i] = true
	}
}

func ComputeUsedResources() {
	var err error

	_, err = ComputeUsedCpus()
	if err != nil {
		log.WithField("err", err).Errorln("Error computing used CPU")
	}

	_, err = ComputeUsedMemory()
	if err != nil {
		log.WithField("err", err).Errorln("Error computing used Memory")
	}
}

func ComputeUsedCpus() (int64, error) {
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

	resources.CPU.Used = cpus

	return cpus, nil
}

func ComputeUsedMemory() (int64, error) {
	var memory int64

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

			memory += cData.Config.Memory
		}
	}

	resources.Memory.Used = memory

	return memory, nil
}

func CheckCoresAvailable(number int) bool {
	defer runtime.Gosched()

	freeCores := 0
	mutex_cpu.RLock()
	for _, free := range resources.CPU.Cores {
		if free {
			freeCores += 1
		}
	}
	mutex_cpu.RUnlock()

	return freeCores >= number
}

func CheckAndSetCores(number int) (string, bool) {
	defer runtime.Gosched()

	cores_int := make([]int, 0, number)
	cores_str := make([]string, 0, number)
	mutex_cpu.Lock()
	for i := 0; i < len(resources.CPU.Cores); i++ {
		if resources.CPU.Cores[i] == true {
			cores_int = append(cores_int, i)
		}

		if len(cores_int) >= number {
			break
		}
	}

	if len(cores_int) < number {
		mutex_cpu.Unlock()
		return "", false
	}

	for _, core := range cores_int {
		resources.CPU.Cores[core] = false
		cores_str = append(cores_str, strconv.Itoa(core))
	}
	mutex_cpu.Unlock()

	return strings.Join(cores_str, ","), true
}

func CheckSpecificCoresAvailable(cpusetcpus string) bool {
	defer runtime.Gosched()

	available := true
	request, err := getCoresNumber(cpusetcpus)
	if err != nil {
		log.WithField("err", err).Errorln("Error Checking available cores")
		return false
	}

	mutex_cpu.RLock()
	for _, req := range request {
		if resources.CPU.Cores[req] == false {
			available = false
		}
	}
	mutex_cpu.RUnlock()

	return available
}

func CheckAndSetSpecificCores(cpusetcpus string) bool {
	defer runtime.Gosched()

	request, err := getCoresNumber(cpusetcpus)
	if err != nil {
		log.WithField("err", err).Errorln("Error Checking available cores")
		return false
	}

	mutex_cpu.Lock()
	// Double loop because I have to assign ALL the requested cores,
	// so I have to check before if all of them are available
	for _, req := range request {
		if resources.CPU.Cores[req] == false {
			mutex_cpu.Unlock()
			return false
		}
	}

	for _, req := range request {
		resources.CPU.Cores[req] = true
	}
	mutex_cpu.Unlock()

	return true
}

func FreeCores(cores string) {
	defer runtime.Gosched()
	toFree, _ := getCoresNumber(cores)
	mutex_cpu.Lock()
	for _, core := range toFree {
		resources.CPU.Cores[core] = true
	}
	mutex_cpu.Unlock()
}

func getCoresNumber(cores string) ([]int, error) {
	cores_str := strings.Split(cores, ",")
	if len(cores_str) > int(resources.CPU.Total) {
		return []int{}, ErrTooManyCores
	}

	cores_int := make([]int, len(cores_str), len(cores_str))
	for i := 0; i < len(cores_str); i++ {
		core_int, err := strconv.Atoi(cores_str[i])
		if err != nil {
			return []int{}, err
		}
		if core_int >= int(resources.CPU.Total) {
			return []int{}, ErrWrongCoreNumber
		}
		cores_int[i] = core_int
	}

	return cores_int, nil
}
