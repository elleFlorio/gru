package autonomic

import (
	log "github.com/Sirupsen/logrus"
)

type Analyze struct{ c_err chan error }

func (p *Analyze) run(stats GruStats) {
	//Analyze stuff
	log.Debugln("I'm analyzing")
}

func computeCpuAvg(service string, stats *GruStats) {

}
