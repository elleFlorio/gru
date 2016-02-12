package resources

import (
	"strconv"
	"strings"
	"testing"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	cfg "github.com/elleFlorio/gru/configuration"
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

func TestCheckAndSetCores(t *testing.T) {
	defer freeCores()

	var assigned string
	var ok bool

	req := 2
	id := "pippo"
	assigned, ok = CheckAndSetCores(req, id)
	assert.True(t, ok)
	assert.Equal(t, "0,1", assigned)
	assert.Equal(t, c_NCORES-req, getAvailableCores())
	_, ok = instanceCores[id]
	assert.True(t, ok)

	freeCores()
	assignSpecificCores([]int{1, 3}, id)
	assigned, ok = CheckAndSetCores(req, id)
	assert.True(t, ok)
	assert.Equal(t, "0,2", assigned)
	assert.Equal(t, 0, getAvailableCores())

	freeCores()
	req_wrong := c_NCORES + 1
	assigned, ok = CheckAndSetCores(req_wrong, id)
	assert.False(t, ok)
	assert.Equal(t, "", assigned)
	assert.Equal(t, c_NCORES, getAvailableCores())

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

	assert.True(t, CheckAndSetSpecificCores(req, id))
	assert.Equal(t, c_NCORES-2, getAvailableCores())

	freeCores()
	assert.False(t, CheckAndSetSpecificCores(req_tooMany, id))
	assert.Equal(t, c_NCORES, getAvailableCores())

	freeCores()
	assert.False(t, CheckAndSetSpecificCores(req_wrongNumber, id))
	assert.Equal(t, c_NCORES, getAvailableCores())

	freeCores()
	assignSpecificCores([]int{1, 3}, id)
	assert.False(t, CheckAndSetSpecificCores(req, id))

}

func TestFreeInstanceCores(t *testing.T) {
	defer freeCores()

	var id string

	id = "pippo"
	assignCores(2, id)
	assert.Equal(t, c_NCORES-2, getAvailableCores())
	assert.True(t, FreeInstanceCores(id))
	assert.Equal(t, c_NCORES, getAvailableCores())

	id = "topolino"
	assert.False(t, FreeInstanceCores(id))

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
