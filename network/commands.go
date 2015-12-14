package network

import (
	"encoding/json"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

const c_COMMAND_ROUTE = "/gru/v1/commands"

type Command struct {
	Name   string
	Target string
}

func SendStartCommand(target string) error {
	cmd := Command{"start", "agent"}
	err := sendCommand(target, cmd)
	if err != nil {
		log.WithField("err", err).Errorln("Error sending command to target ", target)
		return err
	}

	return nil
}

func sendCommand(address string, cmd Command) error {
	var err error
	body, err := json.Marshal(cmd)
	if err != nil {
		log.WithField("err", err).Errorln("Error marshaling command")
	}
	_, err = DoRequest("POST", address+c_COMMAND_ROUTE, body)
	return err
}
