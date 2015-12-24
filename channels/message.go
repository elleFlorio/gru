package channels

import (
	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/enum"
)

type ActionMessage struct {
	Target  *cfg.Service
	Actions enum.Actions
}
