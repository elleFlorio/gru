package autonomic

import (
	log "github.com/Sirupsen/logrus"
)

type Plan struct{}

func (p *Plan) run(analytics GruAnalytics) {
	//Plan stuff
	log.Debugln("I'm planning")
}
