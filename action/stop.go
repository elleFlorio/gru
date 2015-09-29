package action

import (
	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/container"
)

type Stop struct{}

func (p *Stop) Name() string {
	return "stop"
}

func (p *Stop) Initialize() error {
	return nil
}

func (p *Stop) Run(config *GruActionConfig) error {
	err := container.Docker().Client.StopContainer(config.Target, 1)
	if err != nil {
		log.WithFields(log.Fields{
			"id":    config.Target,
			"error": err,
		}).Errorln("Error stopping container")
		return err
	}

	return nil
}
