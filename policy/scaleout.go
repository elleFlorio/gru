package policy

import (
	"github.com/elleFlorio/gru/autonomic/analyzer"
	"github.com/elleFlorio/gru/node"
	"github.com/elleFlorio/gru/service"
)

type ScaleOut struct {
}

func (p *ScaleOut) Name() string {
	return "scaleout"
}

func (p *ScaleOut) Type() string {
	return "proactive"
}

func (p *ScaleOut) Level() string {
	return "service"
}

func (p *ScaleOut) Actions() []string {
	return []string{
		"start",
	}
}

func (p *ScaleOut) Weight(name string, a *analyzer.GruAnalytics) float64 {
	weight := 0.0
	srv, _ := service.GetServiceByName(name)
	cpuMax := srv.Constraints.CpuMax
	maxActive := srv.Constraints.MaxActive
	maxActiveNode := node.GetNodeConfig().Constraints.MaxInstances
	cpuTot := a.Service[name].CpuTot
	curActive := len(a.Service[name].Instances.Active) + len(a.Service[name].Instances.Pending)
	curActiveNode := len(a.System.Instances.Active) + len(a.System.Instances.Pending)

	// check if the constraints of service are not specified
	if cpuMax == 0.0 {
		cpuMax = 1.0 / float64(node.GetNodeConfig().Constraints.MaxInstances)
	}
	if maxActive == 0 {
		maxActive = maxActiveNode
	}

	// update the scaling up threshold according to active instances
	cpuMax = cpuMax * float64(curActive)

	if curActive < maxActive && curActiveNode < maxActiveNode {
		if cpuTot > cpuMax {
			weight = (cpuTot - cpuMax) / (1.0 - cpuMax)
		}
	}

	return weight
}

func (p *ScaleOut) Target() string {
	return "container"
}

func (p *ScaleOut) TargetStatus() string {
	return "stopped"
}
