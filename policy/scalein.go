package policy

import (
	"github.com/elleFlorio/gru/autonomic/analyzer"
	"github.com/elleFlorio/gru/node"
	"github.com/elleFlorio/gru/service"
)

type ScaleIn struct{}

func (p *ScaleIn) Name() string {
	return "scalein"
}

func (p *ScaleIn) Type() string {
	return "proactive"
}

func (p *ScaleIn) Level() string {
	return "service"
}

func (p *ScaleIn) Actions() []string {
	return []string{
		"stop",
	}
}

func (p *ScaleIn) Weight(name string, a *analyzer.GruAnalytics) float64 {
	weight := 0.0
	srv, _ := service.GetServiceByName(name)
	cpuMin := srv.Constraints.CpuMin
	cpuMax := srv.Constraints.CpuMax
	minActive := srv.Constraints.MinActive
	curActive := len(a.Service[name].Instances.Active)

	// check if the constraints of service are not specified
	if cpuMax == 0.0 {
		cpuMax = 1.0 / float64(node.GetNodeConfig().Constraints.MaxInstances)
	}
	if cpuMin == 0 {
		cpuMin = float64(curActive-1) * cpuMax
	}

	if curActive > minActive {
		cpuTot := a.Service[name].CpuTot
		if cpuTot < cpuMin {
			weight = (cpuMin - cpuTot) / cpuMin
		}
	}

	return weight
}

func (p *ScaleIn) Target() string {
	return "container"
}

func (p *ScaleIn) TargetStatus() string {
	return "running"
}
