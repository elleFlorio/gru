package action

import (
	"errors"

	"github.com/elleFlorio/gru/enum"
)

type ActionExecutor interface {
	Type() enum.Action
	Run(Action) error
}

var (
	actions         []ActionExecutor
	ErrNotSupported = errors.New("action not supported")
)

func init() {
	actions = []ActionExecutor{
		&NoAction{},
		&Start{},
		&Stop{},
		&Remove{},
	}
}

func Get(aType enum.Action) ActionExecutor {
	var act ActionExecutor
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
