package evaluator

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	cfg "github.com/elleFlorio/gru/configuration"
	srv "github.com/elleFlorio/gru/service"
)

func init() {
	srv.SetMockServices()
	setMockExpressions()
}

func TestBuildExpression(t *testing.T) {
	metrics := map[string]float64{
		"M1": 0.5,
		"M2": 0.8,
	}
	constraints := map[string]float64{
		"C1": 0.8,
		"C2": 2.0,
	}

	expr1 := cfg.Expression{
		Expr:        "M1 + M2 / C1",
		Metrics:     []string{"M1", "M2"},
		Constraints: []string{"C1"},
	}
	expct1 := "0.5 + 0.8 / 0.8"

	expr2 := cfg.Expression{
		Expr:        "(M1 + C1) * (M2 + C2)",
		Metrics:     []string{"M1", "M2"},
		Constraints: []string{"C1", "C2"},
	}
	expct2 := "(0.5 + 0.8) * (0.8 + 2)"

	expr3 := cfg.Expression{
		Expr:        "M3 + M2 / C1",
		Metrics:     []string{"M3", "M2"},
		Constraints: []string{"C1"},
	}
	expct3 := "noexp"

	expr4 := cfg.Expression{
		Expr:        "M1 + M2 / C3",
		Metrics:     []string{"M1", "M2"},
		Constraints: []string{"C3"},
	}
	expct4 := "noexp"

	toBuild1 := buildExpression(expr1, metrics, constraints)
	toBuild2 := buildExpression(expr2, metrics, constraints)
	toBuild3 := buildExpression(expr3, metrics, constraints)
	toBuild4 := buildExpression(expr4, metrics, constraints)
	assert.Equal(t, expct1, toBuild1)
	assert.Equal(t, expct2, toBuild2)
	assert.Equal(t, expct3, toBuild3)
	assert.Equal(t, expct4, toBuild4)
}

func TestComputeMetricAnalytics(t *testing.T) {
	metrics := map[string]float64{
		"M1": 0.2,
		"M2": 0.8,
	}
	constraints := map[string]float64{
		"C1": 1.0,
		"C2": 2.0,
	}

	expected1 := map[string]float64{
		"expr1": 0.5,
	}

	expected2 := map[string]float64{
		"expr2": 0.6,
	}

	expected3 := map[string]float64{
		"expr3": 0.0,
	}

	var result map[string]float64

	srv.SetServiceConstraints("service1", constraints)
	result = ComputeMetricAnalytics("service1", metrics)
	assert.Equal(t, expected1, result)

	srv.SetServiceConstraints("service2", constraints)
	result = ComputeMetricAnalytics("service2", metrics)
	assert.Equal(t, expected2, result)

	srv.SetServiceConstraints("service3", constraints)
	result = ComputeMetricAnalytics("service3", metrics)
	assert.Equal(t, expected3, result)

}

func setMockExpressions() {
	expr1 := cfg.Expression{
		Analytic: "expr1",
		Expr:     "(M1 + M2) / 2",
		Metrics:  []string{"M1", "M2"},
	}
	expr2 := cfg.Expression{
		Analytic:    "expr2",
		Expr:        "M1 / C1 + M2 / C2",
		Metrics:     []string{"M1", "M2"},
		Constraints: []string{"C1", "C2"},
	}
	expr3 := cfg.Expression{
		Analytic:    "expr3",
		Expr:        "M1 * C1 + M3 / C3",
		Metrics:     []string{"M1", "M3"},
		Constraints: []string{"C1", "C3"},
	}

	expressions := map[string]cfg.Expression{
		"expr1": expr1,
		"expr2": expr2,
		"expr3": expr3,
	}

	cfg.SetExpr(expressions)
}
