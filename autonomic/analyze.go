package autonomic

import (
	log "github.com/Sirupsen/logrus"
)

type Analyze struct{}

func (p *Analyze) run() {
	//Analyze stuff
	log.Debugln("I'm analyzing")
}
