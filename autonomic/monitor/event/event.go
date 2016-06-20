package event

import "github.com/elleFlorio/gru/enum"

type Event struct {
	Type     string
	Service  string
	Image    string
	Instance string
	Status   enum.Status
}
