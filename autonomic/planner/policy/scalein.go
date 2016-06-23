package policy

import (
	"math"

	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/data"
	"github.com/elleFlorio/gru/enum"
	srv "github.com/elleFlorio/gru/service"
	"github.com/elleFlorio/gru/utils"
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
	if !cfg.GetPolicy().Scalein.Enable {
		return scaleinPolicies
	}

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
	service, _ := srv.GetServiceByName(name)
	inst_run := len(service.Instances.Running)
	inst_pen := len(service.Instances.Pending)

	if inst_run < 1 {
		return 0.0
	}

	baseServices := cfg.GetNodeConstraints().BaseServices
	if (inst_pen+inst_run) <= 1 && utils.ContainsString(baseServices, name) {
		return 0.0
	}

	analytics := srv.GetServiceExpressionsList(name)
	threshold := cfg.GetPolicy().Scalein.Threshold
	weights := []float64{}

	for _, value := range clusterData.Service[name].Data.BaseShared {
		weights = append(weights, p.computeMetricWeight(value, threshold))
	}

	for _, analytic := range analytics {
		value := clusterData.Service[name].Data.UserShared[analytic]
		weights = append(weights, p.computeMetricWeight(value, threshold))
	}

	policyValue := utils.Mean(weights)

	return policyValue
}

func (p *scaleinCreator) computeMetricWeight(value float64, threshold float64) float64 {
	return 1 - (math.Min(value, threshold) / threshold)
}
