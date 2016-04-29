package policy

import (
	"math"

	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/data"
	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/service"
)

type scaleinCreator struct{}

func (p *scaleinCreator) getPolicyName() string {
	return "scalein"
}

func (p *scaleinCreator) listActions() []string {
	return []string{"stop"}
}

func (p *scaleinCreator) createPolicies(srvList []string, clusterData data.Shared) []data.Policy {
	scaleinPolicies := make([]data.Policy, 0, len(srvList))

	for _, name := range srvList {
		policyName := p.getPolicyName()
		policyWeight := p.computeWeight(name, clusterData)
		policyTargets := []string{name}
		policyActions := map[string][]enum.Action{
			name: []enum.Action{enum.STOP, enum.REMOVE},
		}

		scaleinPolicy := data.Policy{
			Name:    policyName,
			Weight:  policyWeight,
			Targets: policyTargets,
			Actions: policyActions,
		}

		scaleinPolicies = append(scaleinPolicies, scaleinPolicy)
	}

	return scaleinPolicies
}

func (p *scaleinCreator) computeWeight(name string, clusterData data.Shared) float64 {
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

	srvShared := clusterData.Service[name]
	// LOAD
	load := srvShared.Load
	thrLoad := cfg.GetTuning().Policy.Scalein.Load
	value_load := math.Min(load, thrLoad)
	weight_load := 1 - value_load/thrLoad
	// CPU
	cpu := srvShared.Cpu
	thrCpu := cfg.GetTuning().Policy.Scalein.Cpu
	value_cpu := math.Min(cpu, thrCpu)
	weight_cpu := 1 - value_cpu/thrCpu
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
