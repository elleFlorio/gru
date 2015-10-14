package action

import (
	"errors"

	"github.com/elleFlorio/gru/enum"
)

type GruAction interface {
	Type() enum.Action
	Run(GruActionConfig) error
}

var (
	actions         []GruAction
	ErrNotSupported = errors.New("action not supported")
)

func init() {
	actions = []GruAction{
		&NoAction{},
		&Start{},
		&Stop{},
	}
}

func Get(aType enum.Action) GruAction {
	var act GruAction
	for _, action := range actions {
		if action.Type() == aType {
			act = action
		}
	}

	return act
}

func List() []enum.Action {
	types := []enum.Action{}

	for _, action := range actions {
		types = append(types, action.Type())
	}

	return types
}
