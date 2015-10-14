package action

import (
	"errors"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/container"
	"github.com/elleFlorio/gru/enum"
)

var ErrNoContainerToStop error = errors.New("No running container to stop")

type Stop struct{}

func (p *Stop) Type() enum.Action {
	return enum.STOP
}

func (p *Stop) Run(config GruActionConfig) error {
	var err error
	running := config.Instances.Running

	if len(running) < 1 {
		log.WithField("error", ErrNoContainerToStop).Errorln("Cannot stop container")
		return ErrNoContainerToStop
	}

	toStop := running[0]
	err = container.Docker().Client.StopContainer(toStop, 1)
	if err != nil {
		log.WithField("error", err).Errorln("Cannot stop container ", toStop)
		return err
	}

	return nil
}
