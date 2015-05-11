package autonomic

import (
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
}
