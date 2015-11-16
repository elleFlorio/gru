package strategy

import (
	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/service"
)

type GruPlan struct {
	Weight  float64
	Target  *service.Service
	Actions []enum.Action
}
