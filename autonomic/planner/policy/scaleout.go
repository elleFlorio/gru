package policy

import (
	"math"

	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/data"
	"github.com/elleFlorio/gru/enum"
	res "github.com/elleFlorio/gru/resources"
	"github.com/elleFlorio/gru/service"
)

type scaleoutCreator struct{}

func (p *scaleoutCreator) getPolicyName() string {
	return "scaleout"
}

func (p *scaleoutCreator) listActions() []string {
	return []string{"start"}
}

func (p *scaleoutCreator) createPolicies(srvList []string, clusterData data.Shared) []data.Policy {
	scaleoutPolicies := make([]data.Policy, 0, len(srvList))

	for _, name := range srvList {
		policyName := p.getPolicyName()
		policyWeight := p.computeWeight(name, clusterData)
		policyTargets := []string{name}
		policyActions := map[string][]enum.Action{
			name: []enum.Action{enum.START},
		}

		scaleoutPolicy := data.Policy{
			Name:    policyName,
			Weight:  policyWeight,
			Targets: policyTargets,
			Actions: policyActions,
		}

		scaleoutPolicies = append(scaleoutPolicies, scaleoutPolicy)
	}

	return scaleoutPolicies
}

func (p *scaleoutCreator) computeWeight(name string, clusterData data.Shared) float64 {
	srv, _ := service.GetServiceByName(name)

	if res.AvailableResourcesService(name) < 1.0 {
		return 0.0
	}

	srvCores := srv.Docker.CpusetCpus
	if srvCores != "" {
		if !res.CheckSpecificCoresAvailable(srvCores) {
			return 0.0
		}
	}

	srvShared := clusterData.Service[name]
	// LOAD
	load := srvShared.Load
	thrLoad := cfg.GetTuning().Policy.Scaleout.Load
	value_load := math.Max(load, thrLoad)
	weight_load := (value_load - thrLoad) / (1 - thrLoad)
	// CPU
	cpu := srvShared.Cpu
	thrCpu := cfg.GetTuning().Policy.Scaleout.Cpu
	value_cpu := math.Max(cpu, thrCpu)
	weight_cpu := (value_cpu - thrCpu) / (1 - thrCpu)
	// MEMORY
	// TODO?

	policyValue := (weight_load + weight_cpu) / 2

	return policyValue
}
