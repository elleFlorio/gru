package autonomic

import (
	log "github.com/Sirupsen/logrus"
)

type Plan struct{}

func (p *Plan) run() {
	//Plan stuff
	log.Debugln("I'm planning")
}
