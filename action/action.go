package action

import (
	"github.com/elleFlorio/gru/service"
)

type GruAction interface {
	Name() string
	Initialize(*service.Service) error
	ComputeWeight() float64
	Execute()
}

var actions []GruAction

func init() {
	actions = []GruAction{
		&ScaleDown{},
		&ScaleUp{},
	}
}

func List() []string {
	names := []string{}

	for _, action := range actions {
		names = append(names, action.Name())
	}

	return names
}
