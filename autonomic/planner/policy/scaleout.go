package policy

import (
	"github.com/elleFlorio/gru/autonomic/analyzer"
	"github.com/elleFlorio/gru/enum"
)

type ScaleOut struct{}

func (p *ScaleOut) Name() string {
	return "scaleout"
}

//TODO find a way to compute a label that make some sense...
func (p *ScaleOut) Label(name string, analytics analyzer.GruAnalytics) enum.Label {

	return enum.WHITE
}

func (p *ScaleOut) Actions() []enum.Action {
	return []enum.Action{
		enum.START,
	}
}
