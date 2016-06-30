package configuration

type AnalyticExpr struct {
	Name        string
	Expr        string
	Metrics     []string
	Constraints []string
}
