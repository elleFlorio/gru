package autonomic

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/samalba/dockerclient"
)

type autoManager struct {
	Docker           *dockerclient.DockerClient
	LoopTimeInterval int
}

var manager autoManager

func NewAutoManager(docker *dockerclient.DockerClient, loopTimeInternval int) *autoManager {
	manager = autoManager{
		docker,
		loopTimeInternval,
	}
	return &manager
}

func (man *autoManager) RunLoop() {
	log.WithField("status", "start").Debugln("Running autonomic loop")

	man.loop()

	log.WithField("status", "done").Debugln("Running autonomic loop")
}

func (man *autoManager) loop() {
	c_err := make(chan error)
	c_stop := make(chan struct{})

	m := Monitor{c_stop, c_err}
	a := Analyze{}
	p := Plan{}
	e := Execute{}

	m.start(man.Docker)

	// Set the ticker for the periodic execution
	ticker := time.NewTicker(time.Duration(man.LoopTimeInterval) * time.Second)

	for {
		select {
		case <-ticker.C:
			stats := m.run()
			a.run(stats)
			p.run()
			e.run()
		case <-c_err:
			log.WithField("status", "error").Debugln("Running autonomic loop")
		case <-c_stop:
			ticker.Stop()
		}
	}
}
