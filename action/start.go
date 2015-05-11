package action

import (
	log "github.com/Sirupsen/logrus"
)

type Start struct{}

func (p *Start) Name() string {
	return "start"
}

func (p *Start) Initialize() error {
	return nil
}

func (p *Start) Run(config *GruActionConfig) error {
	err := config.Client.StartContainer(config.ContainerId, config.HostConf)
	if err != nil {
		log.WithFields(log.Fields{
			"id":    config.ContainerId,
			"error": err,
		}).Errorln("Error starting container")
		return err
	}

	return nil
}
