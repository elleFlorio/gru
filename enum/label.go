package enum

type Label float64

const (
	WHITE  Label = -1
	GREEN  Label = -0.5
	YELLOW Label = 0
	ORANGE Label = 0.5
	RED    Label = 1
)

func GetLabel(value float64) Label {
	switch {
	case value <= 0.2:
		return WHITE
	case value <= 0.4:
		return GREEN
	case value <= 0.6:
		return YELLOW
	case value <= 0.8:
		return ORANGE
	default:
		return RED
	}
}

func (l Label) Value() float64 {
	var v float64
	switch l {
	case WHITE:
		v = -1
	case GREEN:
		v = -0.5
	case YELLOW:
		v = 0
	case ORANGE:
		v = 0.5
	case RED:
		v = 1
	}

	return v
}

func (l Label) ToString() string {
	var s string
	switch l {
	case WHITE:
		s = "WHITE"
	case GREEN:
		s = "GREEN"
	case YELLOW:
		s = "YELLOW"
	case ORANGE:
		s = "ORANGE"
	case RED:
		s = "RED"
	}

	return s
}
