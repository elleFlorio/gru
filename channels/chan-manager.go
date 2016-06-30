package channels

import (
	"errors"

	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/enum"
)

var (
	ch_action        chan ActionMessage
	ch_instances     map[string]chan struct{}
	ch_removal       chan struct{}
	ch_autonomic_err chan error

	needsRemovalNotification bool
)

func init() {
	ch_action = make(chan ActionMessage)
	ch_instances = make(map[string]chan struct{})
	ch_removal = make(chan struct{})
	ch_autonomic_err = make(chan error)

	needsRemovalNotification = false
}

func GetActionChannel() chan ActionMessage {
	return getChannel("action").(chan ActionMessage)
}

func GetRemovalChannel() chan struct{} {
	return getChannel("removal").(chan struct{})
}

func GetAutonomicErrChannel() chan error {
	return getChannel("autonomic_err").(chan error)
}

func getChannel(name string) interface{} {
	switch name {
	case "action":
		return ch_action
	case "removal":
		return ch_removal
	case "autonomic_err":
		return ch_autonomic_err
	}

	return nil
}

func SetRemovalNotification(value bool) {
	needsRemovalNotification = value
}

func NeedsRemovalNotification() bool {
	return needsRemovalNotification
}

func SendActionStartMessage(target *cfg.Service) {
	message := ActionMessage{target, []enum.Action{enum.START}}
	sendActionMessage(message)

}

func SendActionStopMessage(target *cfg.Service) {
	message := ActionMessage{target, []enum.Action{enum.STOP, enum.REMOVE}}
	sendActionMessage(message)

}

func sendActionMessage(message ActionMessage) {
	ch_action <- message
}

func CreateInstanceChannel(id string) chan struct{} {
	ch_instances[id] = make(chan struct{})
	return ch_instances[id]
}

func GetInstanceChannel(id string) (chan struct{}, error) {
	if ch_instance, ok := ch_instances[id]; ok {
		return ch_instance, nil
	}

	return nil, errors.New("Cannot find channel for such instance")
}

func RemoveInstanceChannel(id string) {
	delete(ch_instances, id)
}
