package action

import (
	"errors"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

type GruAction interface {
	Name() string
	Initialize() error
	Run(*GruActionConfig) error
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

func New(name string) (GruAction, error) {
	for _, action := range actions {
		if action.Name() == name {
			log.WithField("name", name).Debugln("Initializing action")
			err := action.Initialize()
			return action, err
		}
	}

	return nil, ErrNotSupported
}

func List() []string {
	names := []string{}

	for _, action := range actions {
		names = append(names, action.Name())
	}

	return names
}
