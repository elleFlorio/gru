package strategy

import (
	"errors"

	log "github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/autonomic/analyzer"
)

type GruStrategy interface {
	Name() string
	Initialize() error
	//TODO
	MakeDecision([]GruPlan, *analyzer.GruAnalytics) (*GruPlan, error)
}

var (
	strategies          []GruStrategy
	ErrorNoSuchStrategy error = errors.New("Strategy not implemented")
	ErrorNoSuchTarget   error = errors.New("Unrecognized Target type")
	ErrorNoSuchStatus   error = errors.New("Unrecognized Target status")
	ErrorNoStoppedCont  error = errors.New("No stopped container for this service")
)

func init() {
	strategies = []GruStrategy{
		&DummyStrategy{},
	}
}

func New(name string) (GruStrategy, error) {

	for _, strategy := range strategies {
		if strategy.Name() == name {
			log.WithField("name", name).Debugln("Initializing strategy")
			err := strategy.Initialize()
			return strategy, err
		}
	}

	return nil, ErrorNoSuchStrategy
}

func List() []string {
	names := []string{}

	for _, strategy := range strategies {
		names = append(names, strategy.Name())
	}

	return names
}
