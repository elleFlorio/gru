package policy

import (
	"github.com/elleFlorio/gru/autonomic/analyzer"
)

type policyCreator interface {
	getPolicyName() string
	createPolicy(analyzer.GruAnalytics, ...string) Policy
	listActions() []string
}

var creators []policyCreator

func init() {
	creators = []policyCreator{
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

func CreatePolicies(srvList []string) []Policy {
	policies := []Policy{}

	return policies
}

func createSwapPairs(srvList []string) map[string][]string {
	pairs := map[string][]string{}

	return pairs
}
