package policy

import (
	"math"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/data"
	"github.com/elleFlorio/gru/enum"
	res "github.com/elleFlorio/gru/resources"
	srv "github.com/elleFlorio/gru/service"
	"github.com/elleFlorio/gru/utils"
)

type swapCreator struct{}

func (p *swapCreator) getPolicyName() string {
	return "swap"
}

func (p *swapCreator) listActions() []string {
	return []string{"stop", "remove", "start"}
}

func (p *swapCreator) createPolicies(srvList []string, clusterData data.Shared) []data.Policy {
	swapPolicies := []data.Policy{}
	if !cfg.GetPolicy().Swap.Enable {
		return swapPolicies
	}

	swapPairs := p.createSwapPairs(srvList)
	for running, inactives := range swapPairs {
		for _, inactive := range inactives {
			policyName := p.getPolicyName()
			policyWeight := p.computeWeight(running, inactive, clusterData)
			policyTargets := []string{running, inactive}
			policyActions := map[string][]enum.Action{
				running:  []enum.Action{enum.STOP, enum.REMOVE},
				inactive: []enum.Action{enum.START},
			}

			swapPolicy := data.Policy{
				Name:    policyName,
				Weight:  policyWeight,
				Targets: policyTargets,
				Actions: policyActions,
			}

			swapPolicies = append(swapPolicies, swapPolicy)
		}
	}

	return swapPolicies
}

func (p *swapCreator) createSwapPairs(srvList []string) map[string][]string {
	pairs := map[string][]string{}

	running := []string{}
	inactive := []string{}

	for _, name := range srvList {
		service, _ := srv.GetServiceByName(name)
		if len(service.Instances.Running) > 0 {
			running = append(running, name)
		} else {
			inactive = append(inactive, name)
		}
	}

	for _, name := range running {
		pairs[name] = inactive
	}

	return pairs
}

func (p *swapCreator) computeWeight(running string, candidate string, clusterData data.Shared) float64 {
	srv_run, _ := srv.GetServiceByName(running)
	srv_cand, _ := srv.GetServiceByName(candidate)
	nRun := len(srv_run.Instances.Running)
	baseServices := cfg.GetNodeConstraints().BaseServices

	if utils.ContainsString(baseServices, running) && nRun < 2 {
		return 0.0
	}

	// If the service has the resources to start without stopping the other
	// there is no reason to swap them
	if res.AvailableResourcesService(candidate) > 0 {
		return 0.0
	}

	// TODO now this works only with homogeneous containers
	// and taking into account only the CPUs. This is not a
	// a good thing, so in the feuture the swap policy should
	// be able to compare the resources needed by each containers
	// and evaulte if it is possible to swap a container with
	// more than one that is active, in order to obtain
	// the requested amount of resources.
	if srv_run.Docker.CPUnumber != srv_cand.Docker.CPUnumber {
		return 0.0
	}

	runShared := clusterData.Service[running]
	candShared := clusterData.Service[candidate]
	policy := cfg.GetPolicy().Swap
	metrics := policy.Metrics
	analytics := policy.Analytics
	threshold := policy.Threshold
	weights := []float64{}

	candValue := 0.0
	runValue := 0.0
	for _, metric := range metrics {
		if value, ok := candShared.Data.BaseShared[metric]; ok {
			candValue = value
		} else {
			candValue = -1.0
		}

		if value, ok := runShared.Data.BaseShared[metric]; ok {
			runValue = value
		} else {
			runValue = -1.0
		}

		if candValue != -1.0 && runValue != -1.0 {
			delta := candValue - runValue
			weight := math.Min(1.0, delta/threshold)
			weights = append(weights, weight)
		} else {
			log.WithFields(log.Fields{
				"metric":    metric,
				"running":   running,
				"candidate": candidate,
			}).Warnln("Cannot compare services: metric not present in both services")
		}
	}

	for _, analytic := range analytics {
		if value, ok := candShared.Data.UserShared[analytic]; ok {
			candValue = value
		} else {
			candValue = -1.0
		}

		if value, ok := runShared.Data.UserShared[analytic]; ok {
			runValue = value
		} else {
			runValue = -1.0
		}

		if candValue != -1.0 && runValue != -1.0 {
			delta := candValue - runValue
			weight := math.Min(1.0, delta/threshold)
			weights = append(weights, weight)
		} else {
			log.WithFields(log.Fields{
				"analytic":  analytic,
				"running":   running,
				"candidate": candidate,
			}).Warnln("Cannot compare services: analytic not present in both services")
		}
	}

	policyValue := math.Max(0.0, utils.Mean(weights))

	return policyValue
}
