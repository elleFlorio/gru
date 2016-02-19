package executor

import (
	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/autonomic/executor/action"
	"github.com/elleFlorio/gru/autonomic/planner/policy"
	ch "github.com/elleFlorio/gru/channels"
	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/service"
)

func ListenToActionMessages() {
	go listen()
}

func listen() {
	ch_action := ch.GetActionChannel()
	for {
		select {
		case msg := <-ch_action:
			log.Debugln("Received action message")
			config := buildConfig(msg.Target)
			executeActions(msg.Actions, config)
		}
	}
}

func Run(chosenPolicy *policy.Policy) {
	log.WithField("status", "init").Debugln("Gru Executor")
	defer log.WithField("status", "done").Debugln("Gru Executor")

	if chosenPolicy == nil {
		log.WithField("err", "No policy to execute").Warnln("Cannot execute actions")
	} else {
		for target, actions := range chosenPolicy.Targets {
			srv := getTargetService(target)
			config := buildConfig(srv)
			executeActions(actions, config)
		}
	}
}

func getTargetService(name string) *cfg.Service {
	var srv *cfg.Service
	srv, err := service.GetServiceByName(name)
	if err != nil {
		srv = &cfg.Service{Name: "noservice"}
	}

	return srv
}

func buildConfig(srv *cfg.Service) action.GruActionConfig {
	actConfig := action.GruActionConfig{}
	actConfig.Service = srv.Name
	actConfig.Instances = srv.Instances
	actConfig.HostConfig = action.CreateHostConfig(srv.Docker)
	actConfig.ContainerConfig = action.CreateContainerConfig(srv.Docker)
	actConfig.ContainerConfig.Image = srv.Image
	actConfig.Parameters.StopTimeout = srv.Docker.StopTimeout

	return actConfig
}

func executeActions(actions []enum.Action, config action.GruActionConfig) {
	var err error
	for _, actionType := range actions {
		act := action.Get(actionType)
		err = act.Run(config)
		if err != nil {
			log.WithFields(log.Fields{
				"err":    err,
				"action": act.Type().ToString(),
			}).Errorln("Action not executed")
		}

		log.WithFields(log.Fields{
			"target": config.Service,
			"action": act.Type().ToString(),
		}).Infoln("Action executed")
	}
}
