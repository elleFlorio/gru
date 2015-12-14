package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/agent"
)

type Command struct {
	Name      string
	Target    string
	Result    string
	Timestamp time.Time
}

func PostCommand(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	cmd, err := readCommand(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	executeCommand(cmd)
	w.WriteHeader(http.StatusAccepted)
}

func readCommand(r *http.Request) (Command, error) {
	var err error
	var cmd Command

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.WithField("err", err).Errorln("Error reading command body")
		return Command{}, err
	}

	if err = r.Body.Close(); err != nil {
		log.WithField("err", err).Errorln("Error closing command body")
		return Command{}, err
	}

	if err = json.Unmarshal(body, &cmd); err != nil {
		log.WithField("err", err).Errorln("Error unmarshaling command body")
		return Command{}, err
	}

	cmd.Timestamp = time.Now()

	log.WithFields(log.Fields{
		"name":      cmd.Name,
		"target":    cmd.Target,
		"timestamp": cmd.Timestamp,
	}).Debugln("Received command")

	return cmd, nil
}

func executeCommand(cmd Command) {
	switch cmd.Name {
	case "start":
		startCommand(cmd)
	default:
		log.Errorln("Unrecognized command name: ", cmd.Name)
	}
}

func startCommand(cmd Command) {
	switch cmd.Target {
	case "agent":
		startAgent()
	default:
		log.WithField("target", cmd.Target).Errorln("Unrecognized target for command start")
	}
}

func startAgent() {
	fmt.Println("Starting agent...")
	log.Infoln("Starting agent...")
	go agent.Run()
}