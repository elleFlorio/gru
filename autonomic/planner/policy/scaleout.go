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
func (p *ScaleOut) Weight(name string, analytics analyzer.GruAnalytics) float64 {
	srv, _ := service.GetServiceByName(name)
	inst_run := len(srv.Instances.Running)
	inst_pen := len(srv.Instances.Pending)

	if (inst_pen + inst_run) > 0 {
		return 0.0
	}

	srvAnalytics := analytics.Service[name]
	if srvAnalytics.Resources.Available < 1.0 {
		return 0.0
	}

	load := srvAnalytics.Load
	cpu := srvAnalytics.Resources.Cpu
	//mem := srvAnalytics.Resources.Memory for now not use memory

	policyValue := (load + cpu) / 2 // I don't know...

	return policyValue
}

func (p *ScaleOut) Actions() []enum.Action {
	return []enum.Action{
		enum.START,
	}
}
