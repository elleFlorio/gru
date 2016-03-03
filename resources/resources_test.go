package resources

import (
	"strconv"
	"strings"
	"testing"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/utils"
)

const c_NCORES = 4
const c_EPSILON = 0.09

func init() {
	setTotalCpus(c_NCORES)
	log.SetLevel(log.ErrorLevel)
}

func setTotalCpus(cores int) {
	for i := 0; i < cores; i++ {
		resources.CPU.Cores[i] = true
	}

	resources.CPU.Total = int64(c_NCORES)
}

func TestGetCoresNumber(t *testing.T) {
	var err error

	cores := "0,1,2,3"
	cores_tooMany := "0,1,2,3,4,5,6"
	cores_wrongNumber := "0,1,85"
	cores_atoi := "1,pippo"

	var n []int
	n, err = getCoresNumber(cores)
	assert.NoError(t, err)
	assert.Equal(t, []int{0, 1, 2, 3}, n)
	n, err = getCoresNumber(cores_tooMany)
	assert.Error(t, err)
	n, err = getCoresNumber(cores_wrongNumber)
	assert.Error(t, err)
	n, err = getCoresNumber(cores_atoi)
	assert.Error(t, err)
}

func TestCheckCoresAvailable(t *testing.T) {
	defer freeCores()

	req := 2
	req_wrong := c_NCORES + 1
	id := "pippo"
	assert.True(t, CheckCoresAvailable(req))
	assert.False(t, CheckCoresAvailable(req_wrong))

	assignCores(3, id)
	assert.False(t, CheckCoresAvailable(req))
}

func TestCheckSpecificCoresAvailable(t *testing.T) {
	defer freeCores()

	req := "0,3"
	req_tooMany := "0,1,2,3,4,5"
	req_wrongNumber := "0,1,85"
	id := "pippo"
	assert.True(t, CheckSpecificCoresAvailable(req))
	assert.False(t, CheckSpecificCoresAvailable(req_tooMany))
	assert.False(t, CheckSpecificCoresAvailable(req_wrongNumber))

	assignCores(3, id)
	assert.False(t, CheckSpecificCoresAvailable(req))
}

func TestGetCoresAvailable(t *testing.T) {
	defer freeCores()

	var assigned string
	var ok bool

	req := 2
	id := "pippo"
	assigned, ok = GetCoresAvailable(req)
	assert.True(t, ok)
	assert.Equal(t, "0,1", assigned)

	freeCores()
	assignSpecificCores([]int{1, 3}, id)
	assigned, ok = GetCoresAvailable(req)
	assert.True(t, ok)
	assert.Equal(t, "0,2", assigned)

	freeCores()
	req_wrong := c_NCORES + 1
	assigned, ok = GetCoresAvailable(req_wrong)
	assert.False(t, ok)
	assert.Equal(t, "", assigned)

}

func TestCheckAndSetSpecificCores(t *testing.T) {
	defer freeCores()

	req := "0,3"
	req_tooMany := "0,1,2,3,4,5"
	req_wrongNumber := "0,1,85"
	id := "pippo"

	assert.NoError(t, CheckAndSetSpecificCores(req, id))
	assert.Equal(t, c_NCORES-2, getAvailableCores())

	freeCores()
	assert.Error(t, CheckAndSetSpecificCores(req_tooMany, id))
	assert.Equal(t, c_NCORES, getAvailableCores())

	freeCores()
	assert.Error(t, CheckAndSetSpecificCores(req_wrongNumber, id))
	assert.Equal(t, c_NCORES, getAvailableCores())

	freeCores()
	assignSpecificCores([]int{1, 3}, id)
	assert.Error(t, CheckAndSetSpecificCores(req, id))

}

func TestFreeInstanceCores(t *testing.T) {
	defer freeCores()

	var id string

	id = "pippo"
	assignCores(2, id)
	assert.Equal(t, c_NCORES-2, getAvailableCores())
	FreeInstanceCores(id)
	assert.Equal(t, c_NCORES, getAvailableCores())

	assignCores(2, id)
	id = "topolino"
	FreeInstanceCores(id)
	assert.Equal(t, c_NCORES-2, getAvailableCores())

}

func assignCores(cores int, id string) {
	assigned := []string{}

	for core, free := range resources.CPU.Cores {
		if free {
			assigned = append(assigned, strconv.Itoa(core))
			resources.CPU.Cores[core] = false
		}

		if len(assigned) == cores {
			break
		}
	}

	instanceCores[id] = strings.Join(assigned, ",")
}

func assignSpecificCores(cores []int, id string) {
	assigned := []string{}
	for _, core := range cores {
		resources.CPU.Cores[core] = false
		assigned = append(assigned, strconv.Itoa(core))
	}

	instanceCores[id] = strings.Join(assigned, ",")
}

func freeCores() {
	for core, _ := range resources.CPU.Cores {
		resources.CPU.Cores[core] = true
	}

	instanceCores = make(map[string]string)
}

func getAvailableCores() int {
	availables := 0
	for _, free := range resources.CPU.Cores {
		if free {
			availables += 1
		}
	}

	return availables
}

func TestAvailableResourcesService(t *testing.T) {
	defer CleanResources()

	name := "test"
	s_over := createService(name, 8, "16G")
	s_bigger := createService(name, 6, "8G")
	s_big := createService(name, 4, "4G")
	s_medium := createService(name, 2, "4G")
	s_low := createService(name, 2, "2G")
	s_lower := createService(name, 1, "1G")
	s_error := createService(name, 1, "error")

	setResources(6, "8G", 6, "8G")
	cfg.SetServices([]cfg.Service{s_over})
	assert.Equal(t, 0.0, AvailableResourcesService(name))
	cfg.SetServices([]cfg.Service{s_bigger})
	assert.Equal(t, 0.0, AvailableResourcesService(name))
	cfg.SetServices([]cfg.Service{s_big})
	assert.Equal(t, 0.0, AvailableResourcesService(name))
	cfg.SetServices([]cfg.Service{s_medium})
	assert.Equal(t, 0.0, AvailableResourcesService(name))
	cfg.SetServices([]cfg.Service{s_low})
	assert.Equal(t, 0.0, AvailableResourcesService(name))
	cfg.SetServices([]cfg.Service{s_lower})
	assert.Equal(t, 0.0, AvailableResourcesService(name))

	setResources(6, "8G", 4, "4G")
	cfg.SetServices([]cfg.Service{s_over})
	assert.Equal(t, 0.0, AvailableResourcesService(name))
	cfg.SetServices([]cfg.Service{s_bigger})
	assert.Equal(t, 0.0, AvailableResourcesService(name))
	cfg.SetServices([]cfg.Service{s_big})
	assert.Equal(t, 0.0, AvailableResourcesService(name))
	cfg.SetServices([]cfg.Service{s_medium})
	assert.Equal(t, 1.0, AvailableResourcesService(name))
	cfg.SetServices([]cfg.Service{s_low})
	assert.Equal(t, 1.0, AvailableResourcesService(name))
	cfg.SetServices([]cfg.Service{s_lower})
	assert.Equal(t, 1.0, AvailableResourcesService(name))

	setResources(6, "8G", 2, "2G")
	cfg.SetServices([]cfg.Service{s_over})
	assert.Equal(t, 0.0, AvailableResourcesService(name))
	cfg.SetServices([]cfg.Service{s_bigger})
	assert.Equal(t, 0.0, AvailableResourcesService(name))
	cfg.SetServices([]cfg.Service{s_big})
	assert.Equal(t, 1.0, AvailableResourcesService(name))
	cfg.SetServices([]cfg.Service{s_medium})
	assert.Equal(t, 1.0, AvailableResourcesService(name))
	cfg.SetServices([]cfg.Service{s_low})
	assert.Equal(t, 1.0, AvailableResourcesService(name))
	cfg.SetServices([]cfg.Service{s_lower})
	assert.Equal(t, 1.0, AvailableResourcesService(name))

	setResources(6, "8G", 0, "0G")
	cfg.SetServices([]cfg.Service{s_error})
	assert.Equal(t, 0.0, AvailableResourcesService(name))
	cfg.SetServices([]cfg.Service{s_over})
	assert.Equal(t, 0.0, AvailableResourcesService(name))
	cfg.SetServices([]cfg.Service{s_bigger})
	assert.Equal(t, 1.0, AvailableResourcesService(name))
	cfg.SetServices([]cfg.Service{s_big})
	assert.Equal(t, 1.0, AvailableResourcesService(name))
	cfg.SetServices([]cfg.Service{s_medium})
	assert.Equal(t, 1.0, AvailableResourcesService(name))
	cfg.SetServices([]cfg.Service{s_low})
	assert.Equal(t, 1.0, AvailableResourcesService(name))
	cfg.SetServices([]cfg.Service{s_lower})
	assert.Equal(t, 1.0, AvailableResourcesService(name))
}

func TestAvailableResources(t *testing.T) {
	setResources(6, "8G", 6, "8G")
	assert.Equal(t, 0.0, AvailableResources())
	setResources(6, "8G", 4, "6G")
	assert.InEpsilon(t, 0.3, AvailableResources(), c_EPSILON)
	setResources(6, "8G", 3, "4G")
	assert.InEpsilon(t, 0.5, AvailableResources(), c_EPSILON)
	setResources(6, "8G", 2, "2G")
	assert.InEpsilon(t, 0.7, AvailableResources(), c_EPSILON)
	setResources(6, "8G", 0, "0G")
	assert.InEpsilon(t, 1.0, AvailableResources(), c_EPSILON)
}

func createService(name string, cpu int, mem string) cfg.Service {
	srvConfig := cfg.ServiceDocker{
		CPUnumber: cpu,
		Memory:    mem,
	}

	srv := cfg.Service{
		Name:   name,
		Docker: srvConfig,
	}

	return srv
}

func TestInitializeServiceAvailablePorts(t *testing.T) {
	defer clearServicePorts()

	service1 := "service1"
	ports1 := map[string]string{
		"50100": "50100-50103",
	}
	InitializeServiceAvailablePorts(service1, ports1)
	servicePorts1 := resources.Network.ServicePorts[service1]
	assert.NotNil(t, servicePorts1)
	status1 := servicePorts1.Status["50100"]
	assert.NotNil(t, status1)
	assert.Empty(t, status1.Occupied)
	assert.Len(t, status1.Available, 4)

	service2 := "service2"
	ports2 := map[string]string{
		"50200": "50200",
	}
	InitializeServiceAvailablePorts(service2, ports2)
	servicePorts2 := resources.Network.ServicePorts[service2]
	assert.NotNil(t, servicePorts2)
	status2 := servicePorts2.Status["50200"]
	assert.NotNil(t, status2)
	assert.Empty(t, status2.Occupied)
	assert.Len(t, status2.Available, 1)

	service3 := "service3"
	ports3 := map[string]string{
		"50300": "pippo",
	}
	InitializeServiceAvailablePorts(service3, ports3)
	servicePorts3 := resources.Network.ServicePorts[service3]
	assert.NotNil(t, servicePorts3)
	status3 := servicePorts3.Status["50300"]
	assert.NotNil(t, status3)
	assert.Empty(t, status3.Occupied)
	assert.Empty(t, status3.Available)
}

func TestRequestPortsForService(t *testing.T) {
	defer clearServicePorts()
	var err error
	createServicePorts()
	service := "pippo"
	port1 := "50100"
	port2 := "50200"

	result, err := RequestPortsForService(service)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Len(t, resources.Network.ServicePorts[service].LastRequested, 2)
	available1 := resources.Network.ServicePorts[service].Status[port1].Available
	occupied1 := resources.Network.ServicePorts[service].Status[port1].Occupied
	assert.Len(t, available1, 4)
	assert.Len(t, occupied1, 0)
	available2 := resources.Network.ServicePorts[service].Status[port2].Available
	occupied2 := resources.Network.ServicePorts[service].Status[port2].Occupied
	assert.Len(t, available2, 2)
	assert.Len(t, occupied2, 0)

	ports := resources.Network.ServicePorts[service]
	status := ports.Status[port1]
	status.Available = []string{}
	ports.Status[port1] = status
	resources.Network.ServicePorts[service] = ports
	_, err = RequestPortsForService(service)
	assert.Error(t, err)
	assert.Empty(t, resources.Network.ServicePorts[service].LastRequested)
}

// func TestAssignPortsToService(t *testing.T) {
// 	defer clearServicePorts()
// 	var err error
// 	createServicePorts()
// 	service := "pippo"
// 	port1 := "50100"
// 	port2 := "50200"

// 	result, err := AssignPortsToService(service)
// 	assert.NoError(t, err)
// 	assert.Len(t, result, 2)
// 	available1 := resources.Network.ServicePorts[service].Status[port1].Available
// 	occupied1 := resources.Network.ServicePorts[service].Status[port1].Occupied
// 	assert.Len(t, available1, 3)
// 	assert.Len(t, occupied1, 1)
// 	available2 := resources.Network.ServicePorts[service].Status[port2].Available
// 	occupied2 := resources.Network.ServicePorts[service].Status[port2].Occupied
// 	assert.Len(t, available2, 1)
// 	assert.Len(t, occupied2, 1)

// 	result, err = AssignPortsToService(service)
// 	assert.NoError(t, err)
// 	assert.Len(t, result, 2)
// 	available1 = resources.Network.ServicePorts[service].Status[port1].Available
// 	occupied1 = resources.Network.ServicePorts[service].Status[port1].Occupied
// 	assert.Len(t, available1, 2)
// 	assert.Len(t, occupied1, 2)
// 	available2 = resources.Network.ServicePorts[service].Status[port2].Available
// 	occupied2 = resources.Network.ServicePorts[service].Status[port2].Occupied
// 	assert.Len(t, available2, 0)
// 	assert.Len(t, occupied2, 2)

// 	result, err = AssignPortsToService(service)
// 	assert.Error(t, err)
// 	assert.Empty(t, result)
// 	available1 = resources.Network.ServicePorts[service].Status[port1].Available
// 	occupied1 = resources.Network.ServicePorts[service].Status[port1].Occupied
// 	assert.Len(t, available1, 2)
// 	assert.Len(t, occupied1, 2)
// 	available2 = resources.Network.ServicePorts[service].Status[port2].Available
// 	occupied2 = resources.Network.ServicePorts[service].Status[port2].Occupied
// 	assert.Len(t, available2, 0)
// 	assert.Len(t, occupied2, 2)
// }

func TestAssignSpecificPortsToService(t *testing.T) {
	defer clearServicePorts()
	var err error
	createServicePorts()
	service := "pippo"
	id := "123456789"
	port1 := "50100"
	port2 := "50200"
	bindings := map[string][]string{
		port1: []string{"50100"},
		port2: []string{"50200"},
	}

	err = AssignSpecifiPortsToService(service, id, bindings)
	assert.NoError(t, err)
	available1 := resources.Network.ServicePorts[service].Status[port1].Available
	occupied1 := resources.Network.ServicePorts[service].Status[port1].Occupied
	assert.Len(t, available1, 3)
	assert.Len(t, occupied1, 1)
	available2 := resources.Network.ServicePorts[service].Status[port2].Available
	occupied2 := resources.Network.ServicePorts[service].Status[port2].Occupied
	assert.Len(t, available2, 1)
	assert.Len(t, occupied2, 1)
	assert.NotEmpty(t, instanceBindings)

	err = AssignSpecifiPortsToService(service, id, bindings)
	assert.Error(t, err)
	available1 = resources.Network.ServicePorts[service].Status[port1].Available
	occupied1 = resources.Network.ServicePorts[service].Status[port1].Occupied
	assert.Len(t, available1, 3)
	assert.Len(t, occupied1, 1)
	available2 = resources.Network.ServicePorts[service].Status[port2].Available
	occupied2 = resources.Network.ServicePorts[service].Status[port2].Occupied
	assert.Len(t, available2, 1)
	assert.Len(t, occupied2, 1)

	bindings[port1] = []string{"50000"}
	bindings[port2] = []string{"50201"}
	err = AssignSpecifiPortsToService(service, id, bindings)
	assert.NoError(t, err)
	available1 = resources.Network.ServicePorts[service].Status[port1].Available
	occupied1 = resources.Network.ServicePorts[service].Status[port1].Occupied
	assert.Len(t, available1, 3)
	assert.Len(t, occupied1, 2)
	available2 = resources.Network.ServicePorts[service].Status[port2].Available
	occupied2 = resources.Network.ServicePorts[service].Status[port2].Occupied
	assert.Len(t, available2, 0)
	assert.Len(t, occupied2, 2)

}

func TestFreePortsFromService(t *testing.T) {
	defer clearServicePorts()
	createServicePorts()
	service := "pippo"
	id := "123456789"
	port1 := "50100"
	assignPort(service, id, port1, []string{"50100"})
	available1 := resources.Network.ServicePorts[service].Status[port1].Available
	occupied1 := resources.Network.ServicePorts[service].Status[port1].Occupied
	assert.Len(t, available1, 3)
	assert.Len(t, occupied1, 1)

	FreePortsFromService(service, id)
	available1 = resources.Network.ServicePorts[service].Status[port1].Available
	occupied1 = resources.Network.ServicePorts[service].Status[port1].Occupied
	assert.Len(t, available1, 4)
	assert.Len(t, occupied1, 0)

	FreePortsFromService(service, id)
	assert.Len(t, available1, 4)
	assert.Len(t, occupied1, 0)

}

func clearServicePorts() {
	resources.Network.ServicePorts = make(map[string]Ports)
}

func createServicePorts() {
	service := "pippo"
	port1 := "50100"
	port2 := "50200"
	portRange1, _ := utils.GetCompleteRange("50100-50103")
	portRange2, _ := utils.GetCompleteRange("50200-50201")
	status1 := PortStatus{
		Available: portRange1,
	}
	status2 := PortStatus{
		Available: portRange2,
	}
	portStatus := map[string]PortStatus{
		port1: status1,
		port2: status2,
	}
	servicePorts := Ports{
		Status:        portStatus,
		LastAssigned:  make(map[string][]string),
		LastRequested: make(map[string]string),
	}
	resources.Network.ServicePorts[service] = servicePorts

}

func assignPort(service string, id string, guest string, host []string) {
	ports := portBindings{
		guest: host,
	}
	instanceBindings[id] = ports
	servicePorts := resources.Network.ServicePorts[service]
	status := servicePorts.Status[guest]
	for _, p := range host {
		for i := 0; i < len(status.Available); i++ {
			if status.Available[i] == p {
				status.Available = append(status.Available[:i], status.Available[i+1:]...)
				status.Occupied = append(status.Occupied, p)
			}
		}
	}
	servicePorts.Status[guest] = status
	resources.Network.ServicePorts[service] = servicePorts
}
