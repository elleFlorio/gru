package policy

import (
	"github.com/elleFlorio/gru/service"
)

type GruPolicy interface {
	Name() string
	Type() string
	Weight(s *service.Service) float64
	Actions() []string
}

var proactive, reactive []GruPolicy

func init() {
	reactive = []GruPolicy{}
	proactive = []GruPolicy{
		&ScaleDown{},
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
