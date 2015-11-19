package policy

import (
	"math"

	"github.com/elleFlorio/gru/autonomic/analyzer"
	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/node"
	"github.com/elleFlorio/gru/service"
)

const c_THRESHOLD_SCALEIN = 0.4

type ScaleIn struct{}

func (p *ScaleIn) Name() string {
	return "scalein"
}

//TODO find a way to compute a label that make some sense...
func (p *ScaleIn) Weight(name string, analytics analyzer.GruAnalytics) float64 {
	srv, _ := service.GetServiceByName(name)
	inst_run := len(srv.Instances.Running)
	inst_pen := len(srv.Instances.Pending)

	if inst_run < 1 {
		return 0.0
	}

	baseServices := node.Config().Constraints.BaseServices
	if (inst_pen+inst_run) <= 1 && contains(baseServices, name) {
		return 0.0
	}

	srvAnalytics := analytics.Service[name]
	// LOAD
	load := srvAnalytics.Load
	value_load := math.Min(load, c_THRESHOLD_SCALEIN)
	weight_load := 1 - value_load/c_THRESHOLD_SCALEIN
	// CPU
	cpu := srvAnalytics.Resources.Cpu
	value_cpu := math.Min(cpu, c_THRESHOLD_SCALEIN)
	weight_cpu := 1 - value_cpu/c_THRESHOLD_SCALEIN
	// MEMORY
	// TODO?

	policyValue := (weight_load + weight_cpu) / 2

	return policyValue
}

func contains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}

	return false
}

func (p *ScaleIn) Actions() []enum.Action {
	return []enum.Action{
		enum.STOP,
	}
}
