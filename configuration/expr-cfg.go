package configuration

type Expression struct {
	Analytic    string
	Expr        string
	Metrics     []string
	Constraints []string
}
