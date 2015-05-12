package autonomic

import (
	log "github.com/Sirupsen/logrus"
)

type Execute struct{}

func (p *Execute) run() {
	//Execute stuff
	log.Debugln("I'm executing")
}
