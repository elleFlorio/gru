package planner

import (
	log "github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/autonomic/analyzer"
)

type planner struct{}

func NewPlanner() *planner {
	return &planner{}
}

func (p *planner) Run(analytics analyzer.GruAnalytics) {
	//Plan stuff
	log.Debugln("I'm planning")
}
