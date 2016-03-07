package action

import (
	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/container"
	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/utils"
)

type Start struct{}

func (p *Start) Type() enum.Action {
	return enum.START
}

func (p *Start) Run(config Action) error {
	var toStart string
	var err error
	paused := config.Instances.Paused
	stopped := config.Instances.Stopped

	if len(paused) > 0 {
		log.WithFields(log.Fields{
			"start":   "Paused",
			"service": config.Service,
		}).Debugln("Starting a paused container")
		toStart = paused[0]
		err = container.Docker().Client.UnpauseContainer(toStart)
		if err != nil {
			return err
		}

		return nil
	}

	if len(stopped) > 0 {
		log.WithFields(log.Fields{
			"start":   "Stopped",
			"service": config.Service,
		}).Debugln("Starting a stopped container")
		toStart = stopped[0]
		err = container.Docker().Client.StartContainer(toStart, config.HostConfig)
		if err != nil {
			return err
		}

		return nil

	}

	log.WithFields(log.Fields{
		"start":   "New",
		"service": config.Service,
	}).Debugln("No stopped/paused container to start: creating new one")
	toStart, err = createNewContainer(config)
	err = container.Docker().Client.StartContainer(toStart, config.HostConfig)
	if err != nil {
		return err
	}

	return nil

}

func createNewContainer(config Action) (string, error) {
	uuid, err := utils.GenerateUUID()
	name := config.Service + "_" + uuid
	id, err := container.Docker().Client.CreateContainer(config.ContainerConfig, name, nil)
	if err != nil {
		log.WithField("err", err).Errorln("Cannot create a new container for service ", config.Service)
		return "", err
	}
	return id, nil
}
