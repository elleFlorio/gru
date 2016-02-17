package policy

import (
	"math"

	"github.com/elleFlorio/gru/autonomic/analyzer"
	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/service"
)

const c_THRESHOLD_SCALEIN_LOAD = 0.3
const c_THRESHOLD_SCALEIN_CPU = 0.3

type scaleinCreator struct{}

func (p *scaleinCreator) getPolicyName() string {
	return "scalein"
}

func (p *scaleinCreator) listActions() []string {
	return []string{"stop"}
}

func (p *scaleinCreator) createPolicies(srvList []string, analytics analyzer.GruAnalytics) []Policy {
	scaleinPolicies := make([]Policy, 0, len(srvList))

	for _, name := range srvList {
		policyName := p.getPolicyName()
		policyWeight := p.computeWeight(name, analytics)
		policyTargets := map[string]enum.Action{
			name: enum.STOP,
		}

		scaleinPolicy := Policy{
			Name:    policyName,
			Weight:  policyWeight,
			Targets: policyTargets,
		}

		scaleinPolicies = append(scaleinPolicies, scaleinPolicy)
	}

	return scaleinPolicies
}

func (p *scaleinCreator) computeWeight(name string, analytics analyzer.GruAnalytics) float64 {
	srv, _ := service.GetServiceByName(name)
	inst_run := len(srv.Instances.Running)
	inst_pen := len(srv.Instances.Pending)

	if inst_run < 1 {
		return 0.0
	}

	baseServices := cfg.GetNodeConstraints().BaseServices
	if (inst_pen+inst_run) <= 1 && p.contains(baseServices, name) {
		return 0.0
	}

	srvAnalytics := analytics.Service[name]
	// LOAD
	load := srvAnalytics.Load
	value_load := math.Min(load, c_THRESHOLD_SCALEIN_LOAD)
	weight_load := 1 - value_load/c_THRESHOLD_SCALEIN_LOAD
	// CPU
	cpu := srvAnalytics.Resources.Cpu
	value_cpu := math.Min(cpu, c_THRESHOLD_SCALEIN_CPU)
	weight_cpu := 1 - value_cpu/c_THRESHOLD_SCALEIN_CPU
	// MEMORY
	// TODO?

	policyValue := (weight_load + weight_cpu) / 2

	return policyValue
}

func (p *scaleinCreator) contains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}

	return false
}
