package network

import (
	"encoding/json"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

const c_COMMAND_ROUTE = "/gru/v1/commands"

type Command struct {
	Name   string
	Target string
	Object interface{}
}

func SendStartCommand(dest string) error {
	cmd := Command{"start", "agent", nil}
	err := sendCommand(dest, cmd)
	if err != nil {
		log.WithField("err", err).Errorln("Error sending command to destination ", dest)
		return err
	}

	return nil
}

func SendUpdateCommand(dest string, target string, obj interface{}) error {
	cmd := Command{"update", target, obj}
	err := sendCommand(dest, cmd)
	if err != nil {
		log.WithField("err", err).Errorln("Error sending command to destination ", dest)
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
