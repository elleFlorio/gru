package policy

import (
	"github.com/elleFlorio/gru/autonomic/analyzer"
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

func (p *ScaleDown) Weight(s *service.Service, a *analyzer.GruAnalytics) float64 {
	weight := 0.0
	cpuAvg := a.Service[s.Name].CpuAvg
	if cpuAvg < s.Constraints.CpuMin {
		weight = (s.Constraints.CpuMin - cpuAvg) / s.Constraints.CpuMin
	}
	return weight
}

func (p *ScaleDown) Target() string {
	return "container"
}
