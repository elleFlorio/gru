package policy

import (
	"github.com/elleFlorio/gru/service"
)

type ScaleUp struct {
}

func (p *ScaleUp) Name() string {
	return "scaleup"
}

func (p *ScaleUp) Type() string {
	return "proactive"
}

func (p *ScaleUp) Actions() []string {
	return []string{
		"start",
	}
}

func (p *ScaleUp) Weight(s *service.Service) float64 {
	weight := 0.0
	if s.CpuAvg > s.Constraints.CpuMax {
		weight = (s.CpuAvg - s.Constraints.CpuMax) / (1.0 - s.Constraints.CpuMax)
	}
	return weight
}
