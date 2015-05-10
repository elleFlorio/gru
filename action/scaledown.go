package action

import (
	"github.com/elleFlorio/gru/service"
)

type ScaleDown struct {
	Service string
	CpuMin  float64
	CpuAvg  float64
	Weight  float64
}

func (p *ScaleDown) Name() string {
	return "scaledown"
}

func (p *ScaleDown) Initialize(service *service.Service) error {
	p.Service = service.Name
	p.CpuMin = service.Constraints.CpuMin
	p.CpuAvg = service.CpuAvg
	p.Weight = 0.0
	return nil
}

func (p *ScaleDown) ComputeWeight() float64 {
	if p.CpuAvg < p.CpuMin {
		p.Weight = (p.CpuMin - p.CpuAvg) / p.CpuMin
	}

	return p.Weight
}

func (p *ScaleDown) Execute() {
}
