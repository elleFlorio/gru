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

	ErrTooManyCores        = errors.New("Number of requested cores exceeds the total number of available ones")
	ErrWrongCoreNumber     = errors.New("The number of the specified core is not correct")
	errCPUsAlreadyOccupied = errors.New("CPUs already occupied")
	ErrNoAvailablePorts    = errors.New("No available ports for service")
	errPortAlreadyOccupied = errors.New("Port already occupied")
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

	mutex_port.Lock()
	resources.Network.ServicePorts = make(map[string]Ports)
	mutex_port.Unlock()
}

func GetResources() *Resource {
	return &resources
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

// ############### CPU ######################

func GetInstanceCores(id string) string {
	if cores, ok := instanceCores[id]; ok {
		return cores
	}

	return ""
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
// func CheckAndSetCores(number int, id string) (string, bool) {
// 	defer runtime.Gosched()

// 	cores_int := make([]int, 0, number)
// 	cores_str := make([]string, 0, number)
// 	mutex_cpu.Lock()
// 	for i := 0; i < len(resources.CPU.Cores); i++ {
// 		if resources.CPU.Cores[i] == true {
// 			cores_int = append(cores_int, i)
// 		}

// 		if len(cores_int) >= number {
// 			break
// 		}
// 	}

// 	if len(cores_int) < number {
// 		log.Errorln("Error assigning cores: number of free cores < ", number)
// 		mutex_cpu.Unlock()
// 		return "", false
// 	}

// 	for _, core := range cores_int {
// 		resources.CPU.Cores[core] = false
// 		cores_str = append(cores_str, strconv.Itoa(core))
// 	}
// 	mutex_cpu.Unlock()

// 	// Record assigned cores to instance
// 	cores := strings.Join(cores_str, ",")
// 	mutex_instance.Lock()
// 	instanceCores[id] = cores
// 	mutex_instance.Unlock()

// 	return cores, true
// }

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

func CheckAndSetSpecificCores(cpusetcpus string, id string) error {
	defer runtime.Gosched()

	request, err := getCoresNumber(cpusetcpus)
	if err != nil {
		log.WithField("err", err).Errorln("Error Checking available cores")
		return err
	}

	mutex_cpu.Lock()
	// Double loop because I have to assign ALL the requested cores,
	// so I have to check before if all of them are available
	for _, req := range request {
		if resources.CPU.Cores[req] == false {
			mutex_cpu.Unlock()
			return errCPUsAlreadyOccupied
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

	return nil
}

func FreeInstanceCores(id string) {
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

	}

	log.Errorln("Error freeing cores: unrecognized instance ", id)

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

// ############## NETWORK ################

func InitializeServiceAvailablePorts(name string, ports map[string]string) {
	defer runtime.Gosched()
	mutex_port.Lock()

	servicePorts := resources.Network.ServicePorts[name]

	for guest, host := range ports {
		servicePorts.LastAssigned = make(map[string]string)
		servicePorts.Status = make(map[string]PortStatus)

		hostRange, err := utils.GetCompleteRange(host)
		if err != nil {
			log.WithFields(log.Fields{
				"err":     err,
				"service": name,
				"guest":   guest,
				"host":    host,
			}).Warnln("Cannot compute host port range for guest port")
		}

		status := PortStatus{
			Available: hostRange,
			Occupied:  []string{},
		}

		servicePorts.Status[guest] = status
	}
	resources.Network.ServicePorts[name] = servicePorts

	mutex_port.Unlock()
}

func AssignPortsToService(name string) (map[string]string, error) {
	defer runtime.Gosched()
	mutex_port.Lock()
	servicePorts := resources.Network.ServicePorts[name]

	for guest, host := range servicePorts.Status {
		if len(host.Available) < 1 {
			servicePorts.LastAssigned = make(map[string]string)
			resources.Network.ServicePorts[name] = servicePorts

			mutex_port.Unlock()
			return servicePorts.LastAssigned, ErrNoAvailablePorts
		}

		status := PortStatus{}
		status.Available, status.Occupied = moveItem(host.Available, host.Occupied)
		servicePorts.Status[guest] = status
		servicePorts.LastAssigned[guest] = status.Occupied[len(status.Occupied)-1]
	}

	resources.Network.ServicePorts[name] = servicePorts

	mutex_port.Unlock()
	return servicePorts.LastAssigned, nil
}

func AssignSpecifiPortsToService(name string, portBindings map[string][]string) error {
	defer runtime.Gosched()
	mutex_port.Lock()
	servicePorts := resources.Network.ServicePorts[name]

	for guest, bindings := range portBindings {
		status := servicePorts.Status[guest]
		for _, binding := range bindings {
			if contains(binding, status.Occupied) {
				mutex_port.Unlock()
				return errPortAlreadyOccupied
			}

			if contains(binding, status.Available) {
				status.Available, status.Occupied = moveSpecificItem(binding, status.Available, status.Occupied)
			} else {
				// TODO Log message
				status.Occupied = append(status.Occupied, binding)
			}
		}

		servicePorts.Status[guest] = status
	}

	resources.Network.ServicePorts[name] = servicePorts

	mutex_port.Unlock()
	return nil
}

func contains(item string, slice []string) bool {
	for _, element := range slice {
		if element == item {
			return true
		}
	}

	return false
}

func FreePortsFromService(name string, portBindings map[string][]string) {
	defer runtime.Gosched()
	mutex_port.Lock()
	servicePorts := resources.Network.ServicePorts[name]

	for guest, bindings := range portBindings {
		status := servicePorts.Status[guest]
		for _, binding := range bindings {
			if contains(binding, status.Occupied) {
				status.Occupied, status.Available = moveSpecificItem(binding, status.Occupied, status.Available)
			} else {
				log.WithFields(log.Fields{
					"service": name,
					"guest":   guest,
					"host":    binding,
				}).Warnln("Cannot find port in occupied list")
			}

		}
		servicePorts.Status[guest] = status
	}

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

func moveSpecificItem(item string, source []string, dest []string) ([]string, []string) {
	if len(source) == 0 {
		return source, dest
	}

	index := 0
	for i := 0; i < len(source); i++ {
		if item == source[i] {
			index = i
		}
	}

	source = append(source[:index], source[index+1:]...)
	dest = append(dest, item)

	return source, dest

}

func GetAssignedPorts(name string) map[string]string {
	return resources.Network.ServicePorts[name].LastAssigned
}
