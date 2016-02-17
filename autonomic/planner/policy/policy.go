package policy

import (
	"github.com/elleFlorio/gru/enum"
)

type Policy struct {
	Name    string
	Weight  float64
	Targets map[string]enum.Action
}
