package policy

import (
	"github.com/elleFlorio/gru/autonomic/analyzer"
)

type GruPolicy interface {
	Name() string
	Type() string
	Level() string
	Weight(name string, a *analyzer.GruAnalytics) float64
	Target() string
	TargetStatus() string
	Actions() []string
}

var proactive, reactive []GruPolicy

func init() {
	reactive = []GruPolicy{}
	proactive = []GruPolicy{
		&ScaleIn{},
		&ScaleOut{},
	}
}

func GetPolicies(pType string) []GruPolicy {
	if pType == "reactive" {
		return reactive
	} else {
		return proactive
	}
}

func List(pType string) []string {
	names := []string{}
	var policies []GruPolicy

	if pType == "reactive" {
		policies = reactive
	} else {
		policies = proactive
	}

	for _, policy := range policies {
		names = append(names, policy.Name())
	}

	return names
}
