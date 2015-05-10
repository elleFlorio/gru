package action

import (
	"github.com/elleFlorio/gru/service"
)

type GruAction interface {
	Name() string
	Initialize() error
}

var actions []GruAction

func init() {
	actions = []GruAction{}
}

func New(name string) (GruAction, error) {

}

func List() []string {
	names := []string{}

	for _, action := range actions {
		names = append(names, action.Name())
	}

	return names
}
