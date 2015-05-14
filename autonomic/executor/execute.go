package executor

import (
	log "github.com/Sirupsen/logrus"
)

type executor struct{}

func NewExecutor() *executor {
	return &executor{}
}

func (p *executor) Run() {
	//Execute stuff
	log.Debugln("I'm executing")
}
