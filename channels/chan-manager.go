package channels

import (
	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/enum"
)

var (
	ch_action chan ActionMessage
)

func init() {
	ch_action = make(chan ActionMessage)
}

func GetActionChannel() chan ActionMessage {
	return getChannel("action").(chan ActionMessage)
}

func getChannel(name string) interface{} {
	switch name {
	case "action":
		return ch_action
	}

	return nil
}

func SendActionStartMessage(target *cfg.Service) {
	message := ActionMessage{target, []enum.Action{enum.START}}
	sendActionMessage(message)

}

func SendActionStopMessage(target *cfg.Service) {
	message := ActionMessage{target, []enum.Action{enum.STOP}}
	sendActionMessage(message)

}

func sendActionMessage(message ActionMessage) {
	ch_action <- message
}
