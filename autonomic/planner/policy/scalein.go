package policy

import (
	"github.com/elleFlorio/gru/autonomic/analyzer"
	"github.com/elleFlorio/gru/enum"
)

type ScaleIn struct{}

func (p *ScaleIn) Name() string {
	return "scalein"
}

//TODO find a way to compute a label that make some sense...
func (p *ScaleIn) Label(name string, analytics analyzer.GruAnalytics) enum.Label {

	return enum.WHITE
}

func (p *ScaleIn) Actions() []enum.Action {
	return []enum.Action{
		enum.STOP,
	}
}
