package resources

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

const c_NCORES = 4

func init() {
	setTotalCpus(c_NCORES)
}

func setTotalCpus(cores int) {
	for i := 0; i < cores; i++ {
		resources.CPU.Cores[i] = true
	}
}

func TestCheckCoresAvailable(t *testing.T) {
	defer freeCores()

	req := 2
	req_wrong := c_NCORES + 1
	assert.True(t, CheckCoresAvailable(req))
	assert.False(t, CheckCoresAvailable(req_wrong))

	assignCores(3)
	assert.False(t, CheckCoresAvailable(req))
}

func TestCheckSpecificCoresAvailable(t *testing.T) {

}

func TestCheckAndSetCores(t *testing.T) {
	defer freeCores()
	var assigned string
	var ok bool

	req := 2
	assigned, ok = CheckAndSetCores(req)
	assert.True(t, ok)
	assert.Equal(t, "0,1", assigned)
	assert.Equal(t, c_NCORES-req, getAvailableCores())

	freeCores()
	assignSpecificCores([]int{1, 3})
	assigned, ok = CheckAndSetCores(req)
	assert.True(t, ok)
	assert.Equal(t, "0,2", assigned)
	assert.Equal(t, 0, getAvailableCores())

	freeCores()
	req_wrong := c_NCORES + 1
	assigned, ok = CheckAndSetCores(req_wrong)
	assert.False(t, ok)
	assert.Equal(t, "", assigned)
	assert.Equal(t, c_NCORES, getAvailableCores())

}

func assignCores(cores int) {
	assigned := 0

	for core, free := range resources.CPU.Cores {
		if free {
			assigned += 1
			resources.CPU.Cores[core] = false
		}

		if assigned == cores {
			break
		}
	}
}

func assignSpecificCores(cores []int) {
	for _, core := range cores {
		resources.CPU.Cores[core] = false
	}

}

func freeCores() {
	for core, _ := range resources.CPU.Cores {
		resources.CPU.Cores[core] = true
	}
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
