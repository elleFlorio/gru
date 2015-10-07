package enum

type Action float64

const (
	NOACTION Action = iota
	START    Action = iota
	STOP     Action = iota
)

func (a Action) Value() float64 {
	var v float64
	switch {
	case a == NOACTION:
		v = 0.0
	case a == START:
		v = 1.0
	case a == STOP:
		v = 2.0
	}

	return v
}

func (a Action) ToString() string {
	var s string
	switch {
	case a == NOACTION:
		s = "NOACTION"
	case a == START:
		s = "START"
	case a == STOP:
		s = "STOP"
	}

	return s
}