package strategy

import (
	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/enum"
)

type GruPlan struct {
	Policy  string
	Weight  float64
	Target  *cfg.Service
	Actions enum.Actions
}
