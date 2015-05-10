package autonomic

import (
	log "github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/action"
)

func plan(a_data *analyzeData) action.Action {
	//Plan stuff
	log.Debugln("I'm planning")
	actions := []action.Action{}
	act := makeDecision(actions)
	return act
}

func makeDecision(actions []action.Action) action.Action {
	//TODO
	return action.Action{}
}
