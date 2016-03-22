package policy

import (
	"github.com/elleFlorio/gru/data"
	"github.com/elleFlorio/gru/enum"
)

type policyCreator interface {
	getPolicyName() string
	createPolicies([]string, data.GruAnalytics) []data.Policy
	listActions() []string
}

var creators []policyCreator

func init() {
	creators = []policyCreator{
		&scaleoutCreator{},
		&scaleinCreator{},
		&swapCreator{},
	}
}

func List() []string {
	names := []string{}
	for _, creator := range creators {
		names = append(names, creator.getPolicyName())
	}

	return names
}

func ListPolicyActions(name string) []string {
	for _, creator := range creators {
		if creator.getPolicyName() == name {
			return creator.listActions()
		}
	}

	return []string{}
}

func CreatePolicies(srvList []string, analytics data.GruAnalytics) []data.Policy {
	policies := []data.Policy{}

	for _, creator := range creators {
		creatorPolicies := creator.createPolicies(srvList, analytics)
		policies = append(policies, creatorPolicies...)
	}

	noaction := createNoActionPolicy(policies)
	policies = append(policies, noaction)

	return policies
}

func createNoActionPolicy(policies []data.Policy) data.Policy {
	max := 0.0
	for _, policy := range policies {
		if policy.Weight > max {
			max = policy.Weight
		}
	}

	policyName := "noaction"
	policyWeight := 1.0 - max
	policyTargets := []string{"noservice"}
	policyActions := map[string][]enum.Action{
		"noservice": []enum.Action{enum.NOACTION},
	}

	noactionPolicy := data.Policy{
		Name:    policyName,
		Weight:  policyWeight,
		Targets: policyTargets,
		Actions: policyActions,
	}

	return noactionPolicy
}
