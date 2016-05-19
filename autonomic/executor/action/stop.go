package action

import (
	"errors"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/container"
	"github.com/elleFlorio/gru/enum"
)

var ErrNoContainerToStop error = errors.New("No active container to stop")

type Stop struct{}

func (p *Stop) Type() enum.Action {
	return enum.STOP
}

func (p *Stop) Run(config Action) error {
	var err error
	var toStop string
	running := config.Instances.Running

	if len(running) < 1 {
		log.WithField("err", ErrNoContainerToStop).Errorln("Cannot stop running container. Trying with pending ones...")

		pending := config.Instances.Pending
		if len(pending) < 1 {
			log.WithField("err", ErrNoContainerToStop).Errorln("Cannot stop pending container")
			return ErrNoContainerToStop
		}

		toStop = pending[0]

	} else {

		toStop = running[0]
	}

	err = container.Docker().Client.StopContainer(toStop, config.Parameters.StopTimeout)
	if err != nil {
		log.WithField("err", err).Errorln("Cannot stop container ", toStop)
		return err
	}

	return nil
}
