package enum

type Status string

const (
	PENDING Status = "pending"
	RUNNING Status = "running"
	STOPPED Status = "stopped"
	PAUSED  Status = "paused"
	UNKNOWN Status = "unknown"
)

func (st Status) Value() float64 {
	var v float64
	switch {
	case st == PENDING:
		v = 0.0
	case st == RUNNING:
		v = 1.0
	case st == STOPPED:
		v = 2.0
	case st == PAUSED:
		v = 3.0
	case st == UNKNOWN:
		v = 4.0
	}

	return v
}

func (st Status) ToString() string {
	var s string
	switch {
	case st == PENDING:
		s = "pending"
	case st == RUNNING:
		s = "running"
	case st == STOPPED:
		s = "stopped"
	case st == PAUSED:
		s = "paused"
	case st == UNKNOWN:
		s = "unknown"
	}

	return s
}
