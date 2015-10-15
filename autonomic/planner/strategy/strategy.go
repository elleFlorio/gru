package strategy

import (
	"errors"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/service"
)

type GruStrategy interface {
	Name() string
	Initialize() error
	MakeDecision([]GruPlan) *GruPlan
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

func MakeDecision(plans []GruPlan) *GruPlan {
	return active().MakeDecision(plans)
}

func createNoActionPlan() GruPlan {
	srv := service.Service{Name: "NoService"}
	return GruPlan{enum.RED, &srv, []enum.Action{enum.NOACTION}}
}
