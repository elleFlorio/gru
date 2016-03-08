package policy

import (
	"math"

	"github.com/elleFlorio/gru/autonomic/analyzer"
	"github.com/elleFlorio/gru/enum"
	res "github.com/elleFlorio/gru/resources"
	"github.com/elleFlorio/gru/service"
)

const c_THRESHOLD_SCALEOUT_LOAD = 0.8
const c_THRESHOLD_SCALEOUT_CPU = 0.8

type scaleoutCreator struct{}

func (p *scaleoutCreator) getPolicyName() string {
	return "scaleout"
}

func (p *scaleoutCreator) listActions() []string {
	return []string{"start"}
}

func (p *scaleoutCreator) createPolicies(srvList []string, analytics analyzer.GruAnalytics) []Policy {
	scaleoutPolicies := make([]Policy, 0, len(srvList))

	for _, name := range srvList {
		policyName := p.getPolicyName()
		policyWeight := p.computeWeight(name, analytics)
		policyTargets := map[string][]enum.Action{
			name: []enum.Action{enum.START},
		}

		scaleoutPolicy := Policy{
			Name:    policyName,
			Weight:  policyWeight,
			Targets: policyTargets,
		}

		scaleoutPolicies = append(scaleoutPolicies, scaleoutPolicy)
	}

	return scaleoutPolicies
}

func (p *scaleoutCreator) computeWeight(name string, analytics analyzer.GruAnalytics) float64 {
	srv, _ := service.GetServiceByName(name)
	inst_run := len(srv.Instances.Running)
	inst_pen := len(srv.Instances.Pending)

	// if (inst_pen + inst_run) > 0 {
	// 	return 0.0
	// }

	srvAnalytics := analytics.Service[name]
	if srvAnalytics.Resources.Available < 1.0 {
		return 0.0
	}

	srvCores := srv.Docker.CpusetCpus
	if srvCores != "" {
		if !res.CheckSpecificCoresAvailable(srvCores) {
			return 0.0
		}
	}

	// LOAD
	load := srvAnalytics.Load
	value_load := math.Max(load, c_THRESHOLD_SCALEOUT_LOAD)
	weight_load := (value_load - c_THRESHOLD_SCALEOUT_LOAD) / (1 - c_THRESHOLD_SCALEOUT_LOAD)
	// CPU
	cpu := srvAnalytics.Resources.Cpu
	value_cpu := math.Max(cpu, c_THRESHOLD_SCALEOUT_CPU)
	weight_cpu := (value_cpu - c_THRESHOLD_SCALEOUT_CPU) / (1 - c_THRESHOLD_SCALEOUT_CPU)
	// MEMORY
	// TODO?

	policyValue := (weight_load + weight_cpu) / 2

	return policyValue
}
