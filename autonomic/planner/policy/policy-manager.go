package policy

import (
	"github.com/elleFlorio/gru/autonomic/analyzer"
)

type policyCreator interface {
	getPolicyName() string
	createPolicies([]string, analyzer.GruAnalytics) []Policy
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

func CreatePolicies(srvList []string, analytics analyzer.GruAnalytics) []Policy {
	policies := []Policy{}

	for _, creator := range creators {
		creatorPolicies := creator.createPolicies(srvList, analytics)
		policies = append(policies, creatorPolicies...)
	}

	return policies
}
