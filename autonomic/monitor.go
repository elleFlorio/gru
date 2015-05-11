package autonomic

import (
	log "github.com/Sirupsen/logrus"
	"github.com/samalba/dockerclient"
)

func statCallBack(id string, stats *dockerclient.Stats, ec chan error, args ...interface{}) {
	//In the callback I should update my stats
}

// Events are: create, destroy, die, exec_create, exec_start, export, kill, oom, pause, restart, start, stop, unpause
func eventCallback(event *dockerclient.Event, ec chan error, args ...interface{}) {
	log.WithFields(log.Fields{
		"from":   event.From,
		"status": event.Status,
	}).Debug("Received event")

}

func monitor() {
	// Monitor stuff
	log.Debug("I'm monitoring...")
}
