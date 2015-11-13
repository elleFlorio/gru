package policy

import (
	"github.com/elleFlorio/gru/autonomic/analyzer"
	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/service"
)

type ScaleIn struct{}

func (p *ScaleIn) Name() string {
	return "scalein"
}

//TODO find a way to compute a label that make some sense...
func (p *ScaleIn) Label(name string, analytics analyzer.GruAnalytics) enum.Label {
	srv, _ := service.GetServiceByName(name)
	inst_run := len(srv.Instances.Running)
	inst_pen := len(srv.Instances.Pending)

	if srv.Constraints.MinActive >= (inst_run + inst_pen) {
		return enum.WHITE
	}

	srvAnalytics := analytics.Service[name]
	load := srvAnalytics.Load.Value()
	cpu := srvAnalytics.Resources.Cpu.Value()
	//mem := srvAnalytics.Resources.Memory.Value() for now not use memory

	policyValue := -(load + cpu) / 2

	return enum.FromLabelValue(policyValue)
}

func (p *ScaleIn) Actions() []enum.Action {
	return []enum.Action{
		enum.STOP,
	}
}
