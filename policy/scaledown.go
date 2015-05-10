package policy

import (
	"github.com/elleFlorio/gru/service"
)

type ScaleDown struct{}

func (p *ScaleDown) Name() string {
	return "scaledown"
}

func (p *ScaleDown) Type() string {
	return "proactive"
}

func (p *ScaleDown) Actions() []string {
	return []string{
		"stop",
	}
}

func (p *ScaleDown) Weight(s *service.Service) float64 {
	weight := 0.0
	if s.CpuAvg < s.Constraints.CpuMin {
		weight = (s.Constraints.CpuMin - s.CpuAvg) / s.Constraints.CpuMin
	}
	return weight
}
