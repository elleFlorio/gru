package enum

type Datatype float64

const (
	STATS     Datatype = iota
	ANALYTICS Datatype = iota
	POLICIES  Datatype = iota
	INFO      Datatype = iota
)

func (d Datatype) Value() float64 {
	var v float64
	switch {
	case d == STATS:
		v = 0.0
	case d == ANALYTICS:
		v = 1.0
	case d == POLICIES:
		v = 2.0
	case d == INFO:
		v = 3.0
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
	case d == POLICIES:
		s = "POLICIES"
	case d == INFO:
		s = "INFO"
	}

	return s
}
