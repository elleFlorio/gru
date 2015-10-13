package enum

type Label float64

const (
	WHITE  Label = -1
	GREEN  Label = -0.5
	YELLOW Label = 0
	ORANGE Label = 0.5
	RED    Label = 1
)

func FromValue(value float64) Label {
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

func ValueFrom(label Label) float64 {
	switch {
	case label == WHITE:
		return 0
	case label == GREEN:
		return 0.2
	case label == YELLOW:
		return 0.4
	case label == ORANGE:
		return 0.6
	default:
		return 0.8
	}
}

func FromLabelValue(l_value float64) Label {
	switch {
	case l_value < -0.5:
		return WHITE
	case l_value >= -0.5 && l_value < 0.0:
		return GREEN
	case l_value >= 0.0 && l_value < 0.5:
		return YELLOW
	case l_value >= 0.5 && l_value < 1:
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
