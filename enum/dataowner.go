package enum

type DataOwner float64

const (
	LOCAL   DataOwner = iota
	CLUSTER DataOwner = iota
)

func (o DataOwner) Value() float64 {
	var v float64
	switch {
	case o == LOCAL:
		v = 0.0
	case o == CLUSTER:
		v = 1.0
	}

	return v
}

func (o DataOwner) ToString() string {
	var s string
	switch {
	case o == LOCAL:
		s = "LOCAL"
	case o == CLUSTER:
		s = "CLUSTER"
	}

	return s
}
