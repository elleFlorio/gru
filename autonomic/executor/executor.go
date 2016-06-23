package executor

import (
	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/autonomic/executor/action"
	ch "github.com/elleFlorio/gru/channels"
	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/data"
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
			executeActions(msg.Target, msg.Actions)
		}
	}
}

func Run(chosenPolicy *data.Policy) {
	log.WithField("status", "init").Debugln("Gru Executor")
	defer log.WithField("status", "done").Debugln("Gru Executor")

	if chosenPolicy == nil {
		log.Warnln("No policy to execute")
		return
	}

	for _, target := range chosenPolicy.Targets {
		actions := chosenPolicy.Actions[target]
		srv := getTargetService(target)
		executeActions(srv, actions)
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

func executeActions(target *cfg.Service, actions []enum.Action) {
	var err error
	for _, actionType := range actions {
		config := buildConfig(target, actionType)
		actExecutor := action.Get(actionType)
		err = actExecutor.Run(config)
		if err != nil {
			log.WithFields(log.Fields{
				"err":    err,
				"action": actionType.ToString(),
			}).Errorln("Action not executed")
		}

		log.WithFields(log.Fields{
			"target": target.Name,
			"action": actionType.ToString(),
		}).Infoln("Action executed")
	}
}
