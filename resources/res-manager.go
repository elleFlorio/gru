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
	"github.com/elleFlorio/gru/utils"
)

var (
	resources     Resource
	instanceCores map[string]string

	mutex_cpu      = sync.RWMutex{}
	mutex_instance = sync.RWMutex{}
	mutex_port     = sync.RWMutex{}

	ErrTooManyCores     = errors.New("Number of requested cores exceeds the total number of available ones")
	ErrWrongCoreNumber  = errors.New("The number of the specified core is not correct")
	ErrNoAvailablePorts = errors.New("No available ports for service")
)

func init() {
	resources = Resource{}
	resources.CPU.Cores = make(map[int]bool)
	resources.Network.ServicePorts = make(map[string]Ports)

	instanceCores = make(map[string]string)
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

func CleanResources() {
	defer runtime.Gosched()

	resources.CPU.Used = 0
	resources.Memory.Used = 0

	mutex_cpu.Lock()
	resources.CPU.Cores = make(map[int]bool)
	mutex_cpu.Unlock()

	mutex_instance.Lock()
	instanceCores = make(map[string]string)
	mutex_instance.Unlock()
}

func GetResources() *Resource {
	return &resources
}

func GetInstanceCores(id string) string {
	if cores, ok := instanceCores[id]; ok {
		return cores
	}

	return ""
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

func AvailableResourcesCPU() float64 {
	return 1.0 - (float64(resources.CPU.Used) / float64(resources.CPU.Total))
}

func AvailableResourcesMemory() float64 {
	return 1.0 - (float64(resources.Memory.Used) / float64(resources.Memory.Total))
}

func AvailableResources() float64 {
	return (AvailableResourcesCPU() + AvailableResourcesMemory()) / 2
}

func AvailableResourcesService(name string) float64 {
	var err error

	nodeCpu := resources.CPU.Total
	nodeUsedCpu := resources.CPU.Used
	nodeMem := resources.Memory.Total
	nodeUsedMem := resources.Memory.Used

	srv, _ := service.GetServiceByName(name)
	srvCpu := srv.Docker.CPUnumber
	log.WithFields(log.Fields{
		"service": name,
		"cpus":    srvCpu,
	}).Debugln("Service cpu resources")

	var srvMem int64
	if srv.Docker.Memory != "" {
		srvMem, err = utils.RAMInBytes(srv.Docker.Memory)
		if err != nil {
			log.WithField("err", err).Warnln("Cannot convert service RAM in Bytes.")
			return 0.0
		}
	} else {
		srvMem = 0
	}

	if nodeCpu < int64(srvCpu) || nodeMem < int64(srvMem) {
		return 0.0
	}

	if (nodeCpu-nodeUsedCpu) < int64(srvCpu) || (nodeMem-nodeUsedMem) < int64(srvMem) {
		return 0.0
	}

	return 1.0
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

func GetCoresAvailable(number int) (string, bool) {
	defer runtime.Gosched()

	cores_str := make([]string, 0, number)
	mutex_cpu.RLock()
	for i := 0; i < len(resources.CPU.Cores); i++ {
		if resources.CPU.Cores[i] == true {
			cores_str = append(cores_str, strconv.Itoa(i))
		}

		if len(cores_str) >= number {
			break
		}
	}

	if len(cores_str) < number {
		log.Errorln("Error getting available cores: number of free cores < ", number)
		mutex_cpu.RUnlock()
		return "", false
	}

	mutex_cpu.RUnlock()

	cores := strings.Join(cores_str, ",")
	return cores, true
}

//DEPRECATED
func CheckAndSetCores(number int, id string) (string, bool) {
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
		log.Errorln("Error assigning cores: number of free cores < ", number)
		mutex_cpu.Unlock()
		return "", false
	}

	for _, core := range cores_int {
		resources.CPU.Cores[core] = false
		cores_str = append(cores_str, strconv.Itoa(core))
	}
	mutex_cpu.Unlock()

	// Record assigned cores to instance
	cores := strings.Join(cores_str, ",")
	mutex_instance.Lock()
	instanceCores[id] = cores
	mutex_instance.Unlock()

	return cores, true
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

func CheckAndSetSpecificCores(cpusetcpus string, id string) bool {
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
		resources.CPU.Cores[req] = false
	}
	mutex_cpu.Unlock()

	// Record assigned cores to instance
	mutex_instance.Lock()
	instanceCores[id] = cpusetcpus
	mutex_instance.Unlock()

	log.WithFields(log.Fields{
		"id":    id,
		"cores": cpusetcpus,
	}).Debugln("Assigned cores to instance")

	return true
}

func FreeInstanceCores(id string) bool {
	defer runtime.Gosched()
	if cores, ok := instanceCores[id]; ok {
		toFree, _ := getCoresNumber(cores)
		mutex_cpu.Lock()
		for _, core := range toFree {
			resources.CPU.Cores[core] = true
		}
		mutex_cpu.Unlock()

		mutex_instance.Lock()
		delete(instanceCores, id)
		mutex_instance.Unlock()

		log.WithFields(log.Fields{
			"id":    id,
			"cores": cores,
		}).Debugln("Released cores of instance")

		return true
	}

	log.Errorln("Error freeing cores: unrecognized instance ", id)

	return false
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

func setServiceAvailablePorts(name string, ports []string) {
	defer runtime.Gosched()
	mutex_port.Lock()
	servicePorts := resources.Network.ServicePorts[name]
	servicePorts.Available = ports
	resources.Network.ServicePorts[name] = servicePorts
	mutex_port.Unlock()
}

func getServiceAvailablePorts(name string) []string {
	defer runtime.Gosched()
	mutex_port.RLock()
	available := resources.Network.ServicePorts[name].Available
	mutex_port.RUnlock()
	return available
}

func AssignPortToService(name string) error {
	defer runtime.Gosched()
	mutex_port.Lock()
	servicePorts := resources.Network.ServicePorts[name]
	available := servicePorts.Available
	occupied := servicePorts.Occupied
	if len(available) < 1 {
		mutex_port.Unlock()
		return ErrNoAvailablePorts
	}
	available, occupied = moveItem(available, occupied)
	servicePorts.LastAssigned = occupied[len(occupied)-1]
	servicePorts.Available = available
	servicePorts.Occupied = occupied
	resources.Network.ServicePorts[name] = servicePorts
	mutex_port.Unlock()

	return nil
}

func FreePortFromService(name string) {
	defer runtime.Gosched()
	mutex_port.Lock()
	servicePorts := resources.Network.ServicePorts[name]
	available := servicePorts.Available
	occupied := servicePorts.Occupied
	occupied, available = moveItem(occupied, available)
	servicePorts.Available = available
	servicePorts.Occupied = occupied
	resources.Network.ServicePorts[name] = servicePorts
	mutex_port.Unlock()
}

func moveItem(source []string, dest []string) ([]string, []string) {
	if len(source) == 0 {
		return source, dest
	}

	index := len(source) - 1
	item := source[index]
	if index == 0 {
		source = []string{}
	} else {
		source = source[0 : index-1]
	}
	dest = append(dest, item)

	return source, dest
}

func GetAssignedPort(name string) string {
	return resources.Network.ServicePorts[name].LastAssigned
}
