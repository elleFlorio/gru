package strategy

import (
	"errors"

	"github.com/elleFlorio/gru/autonomic/planner/policy"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

type GruStrategy interface {
	Name() string
	Initialize() error
	MakeDecision([]policy.Policy) *policy.Policy
}

var (
	strategy            int
	strategies          []GruStrategy
	ErrorNoSuchStrategy error = errors.New("Strategy not implemented")
)

func init() {
	strategies = []GruStrategy{
		&dummyStrategy{},
		&probabilisticStrategy{},
	}
}

func New(name string) (GruStrategy, error) {
	strategy = 0
	for index, strtg := range strategies {
		if strtg.Name() == name {
			strategy = index
			log.WithField("name", name).Debugln("Initializing strategy")
			err := strategies[strategy].Initialize()
			return strategies[strategy], err
		}
	}

	return strategies[strategy], ErrorNoSuchStrategy
}

func List() []string {
	names := []string{}

	for _, strategy := range strategies {
		names = append(names, strategy.Name())
	}

	return names
}

func active() GruStrategy {
	return strategies[strategy]
}

func Name() string {
	return active().Name()
}

func Initialize() error {
	return active().Initialize()
}

func MakeDecision(policies []policy.Policy) *policy.Policy {
	return active().MakeDecision(policies)
}
