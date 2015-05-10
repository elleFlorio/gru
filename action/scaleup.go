package action

import (
	"github.com/elleFlorio/gru/service"
)

type ScaleUp struct {
	Service string
	CpuMax  float64
	CpuAvg  float64
	Weight  float64
}

func (p *ScaleUp) Name() string {
	return "scaleup"
}

func (p *ScaleUp) Initialize(service *service.Service) error {
	p.Service = service.Name
	p.CpuAvg = service.CpuAvg
	p.CpuMax = service.Constraints.CpuMax
	p.Weight = 0.0
	return nil
}

func (p *ScaleUp) ComputeWeight() float64 {
	if p.CpuAvg > p.CpuMax {
		p.Weight = (p.CpuAvg - p.CpuMax) / (1.0 - p.CpuMax)
	}

	return p.Weight
}

func (p *ScaleUp) Execute() {
}
