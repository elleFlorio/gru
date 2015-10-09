package enum

type Datatype float64

const (
	STATS     Datatype = iota
	ANALYTICS Datatype = iota
	PLANS     Datatype = iota
)

func (d Datatype) Value() float64 {
	var v float64
	switch {
	case d == STATS:
		v = 0.0
	case d == ANALYTICS:
		v = 1.0
	case d == PLANS:
		v = 2.0
	}

	return v
}

func (d Datatype) ToString() string {
	var s string
	switch {
	case d == STATS:
		s = "STATS"
	case d == ANALYTICS:
		s = "ANALYTICS"
	case d == PLANS:
		s = "PLANS"
	}

	return s
}
