package action

import (
	log "github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/utils"
)

type Start struct{}

func (p *Start) Name() string {
	return "start"
}

func (p *Start) Initialize() error {
	return nil
}

func (p *Start) Run(config *GruActionConfig) error {
	var err error = nil
	var id string

	// If my target type is a container I have to start a stopped one (the target)
	// Otherwise I have to create a new one starting from its image, then start it
	if config.TargetType == "container" {
		id = config.Target

	} else {
		id, err = createNewContainer(config)

		if err != nil {
			log.WithFields(log.Fields{
				"id":    id,
				"error": err,
			}).Errorln("Error creating container")
		}
	}

	config.Client.StartContainer(id, config.HostConfig)
	if err != nil {
		log.WithFields(log.Fields{
			"id":    config.Target,
			"error": err,
		}).Errorln("Error starting container")
		return err
	}

	return nil
}

func createNewContainer(config *GruActionConfig) (string, error) {
	uuid, err := utils.GenerateUUID()
	name := config.Service + "_" + uuid

	config.ContainerConfig.Image = config.Target
	id, err := config.Client.CreateContainer(config.ContainerConfig, name)

	if err != nil {
		log.WithFields(log.Fields{
			"id":    id,
			"error": err,
		}).Errorln("Error creating container")
		return "", err
	}

	return id, nil
}
