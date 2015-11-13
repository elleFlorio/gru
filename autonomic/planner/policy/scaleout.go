package policy

import (
	"github.com/elleFlorio/gru/autonomic/analyzer"
	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/service"
)

type ScaleOut struct{}

func (p *ScaleOut) Name() string {
	return "scaleout"
}

//TODO find a way to compute a label that make some sense...
func (p *ScaleOut) Label(name string, analytics analyzer.GruAnalytics) enum.Label {
	srv, _ := service.GetServiceByName(name)
	inst_run := len(srv.Instances.Running)
	inst_pen := len(srv.Instances.Pending)

	if srv.Constraints.MaxActive > 0 &&
		(inst_pen+inst_run) >= srv.Constraints.MaxActive {
		return enum.WHITE
	}

	srvAnalytics := analytics.Service[name]
	if srvAnalytics.Resources.Available.Value() > enum.ORANGE.Value() {
		return enum.WHITE
	}

	load := srvAnalytics.Load.Value()
	cpu := srvAnalytics.Resources.Cpu.Value()
	//mem := srvAnalytics.Resources.Memory.Value() for now not use memory
	//resources := srvAnalytics.Resources.Available.Value()

	policyValue := (load + cpu) / 2 // I don't know...

	return enum.FromLabelValue(policyValue)
}

func (p *ScaleOut) Actions() []enum.Action {
	return []enum.Action{
		enum.START,
	}
}
