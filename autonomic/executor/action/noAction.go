package action

import (
	"github.com/elleFlorio/gru/enum"
)

type NoAction struct{}

func (p *NoAction) Type() enum.Action {
	return enum.NOACTION
}

func (p *NoAction) Run(config Action) error {
	return nil
}
