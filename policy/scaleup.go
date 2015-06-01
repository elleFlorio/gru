package policy

import (
	"github.com/elleFlorio/gru/autonomic/analyzer"
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

func (p *ScaleUp) Weight(s *service.Service, a *analyzer.GruAnalytics) float64 {
	weight := 0.0
	cpuAvg := a.Service[s.Name].CpuAvg
	if cpuAvg > s.Constraints.CpuMax {
		weight = (cpuAvg - s.Constraints.CpuMax) / (1.0 - s.Constraints.CpuMax)
	}
	return weight
}

func (p *ScaleUp) Target() string {
	return "image"
}
