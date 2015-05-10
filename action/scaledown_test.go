package action

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elleFlorio/gru/service"
)

func TestComputeWeight(t *testing.T) {

}

func createMockServices() []service.Service {
	mockService1 := service.Service{
		Name:   "service1",
		CpuAvg: 0.8,
		Constraints: service.Constraints{
			CpuMin: 0.3,
		},
	}

	mockService2 := service.Service{
		Name:   "service2",
		CpuAvg: 0.2,
		Constraints: service.Constraints{
			CpuMin: 0.3,
		},
	}

	return []service.Service{mockService1, mockService2}
}
