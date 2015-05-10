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

var reactive []GruPolicy
var proactive []GruPolicy

func init() {
	reactive = []GruPolicy{}
	proactive = []GruPolicy{
		&ScaleDown{},
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
