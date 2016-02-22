package metric

import (
	"errors"
)

var ErrNoService error = errors.New("No metric service")

type noService struct{}

func (ns *noService) Name() string {
	return "noservice"
}

func (ns *noService) Initialize(config map[string]interface{}) error {
	return ErrNoService
}

func (ns *noService) StoreMetrics(metrics GruMetric) error {
	return ErrNoService
}
