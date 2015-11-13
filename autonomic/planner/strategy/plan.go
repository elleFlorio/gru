package strategy

import (
	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/service"
)

type GruPlan struct {
	Label   enum.Label
	Target  *service.Service
	Actions []enum.Action
}
