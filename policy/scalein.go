package policy

import (
	"github.com/elleFlorio/gru/autonomic/analyzer"
	"github.com/elleFlorio/gru/service"
)

type ScaleIn struct{}

func (p *ScaleIn) Name() string {
	return "scalein"
}

func (p *ScaleIn) Type() string {
	return "proactive"
}

func (p *ScaleIn) Actions() []string {
	return []string{
		"stop",
	}
}

// FIXME I cannot use the avg Cpu of the service to scale the container.
// I need to find another metric...
func (p *ScaleIn) Weight(s *service.Service, a *analyzer.GruAnalytics) float64 {
	weight := 0.0

	// TODO
	/*
		minActive := s.Constraints.MinActive
		curActive := len(a.Service[s.Name].Instances.Active) + len(a.Service[s.Name].Instances.Pending)

		if curActive > minActive {
			cpuAvg := a.Service[s.Name].CpuAvg
			if cpuAvg < s.Constraints.CpuMin {
				weight = (s.Constraints.CpuMin - cpuAvg) / s.Constraints.CpuMin
			}
		}
	*/

	return weight
}

func (p *ScaleIn) Target() string {
	return "container"
}

func (p *ScaleIn) TargetStatus() string {
	return "running"
}
