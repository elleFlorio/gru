package analyzer

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/data"
	"github.com/elleFlorio/gru/enum"
	// res "github.com/elleFlorio/gru/resources"
	srv "github.com/elleFlorio/gru/service"
	"github.com/elleFlorio/gru/storage"
	// "github.com/elleFlorio/gru/utils"
)

func init() {
	srv.SetMockServices()
	setMockExpressions()
	storage.New("internal")
	constraints := map[string]float64{
		"C1": 1.0,
		"C2": 2.0,
	}
	for _, service := range srv.List() {
		srv.SetServiceConstraints(service, constraints)
	}
}

func TestComputeServicesAnalytics(t *testing.T) {
	stats := data.CreateMockStats()
	expected := make(map[string]data.AnalyticData)
	expected["service1"] = data.AnalyticData{
		BaseAnalytics: map[string]float64{
			enum.METRIC_CPU_AVG.ToString(): 0.6,
			enum.METRIC_MEM_AVG.ToString(): 0.3,
		},
		UserAnalytics: map[string]float64{
			"expr1": 0.5,
		},
	}
	expected["service2"] = data.AnalyticData{
		BaseAnalytics: map[string]float64{
			enum.METRIC_CPU_AVG.ToString(): 0.1,
			enum.METRIC_MEM_AVG.ToString(): 0.1,
		},
		UserAnalytics: map[string]float64{
			"expr2": 0.6,
		},
	}
	expected["service3"] = data.AnalyticData{
		BaseAnalytics: map[string]float64{
			enum.METRIC_CPU_AVG.ToString(): 0.9,
			enum.METRIC_MEM_AVG.ToString(): 0.8,
		},
		UserAnalytics: map[string]float64{
			"expr3": 0.0,
		},
	}

	analytics := computeServicesAnalytics(stats.Metrics.Service)
	assert.Equal(t, expected, analytics)

}

func TestComputeSystemAnalytics(t *testing.T) {
	stats := data.CreateMockStats()
	expected := data.AnalyticData{
		BaseAnalytics: map[string]float64{
			enum.METRIC_CPU_AVG.ToString(): 0.5,
			enum.METRIC_MEM_AVG.ToString(): 0.4,
		},
	}

	analytics := computeSystemAnalytics(stats.Metrics.System)
	assert.Equal(t, expected, analytics)
}

func setMockExpressions() {
	expr1 := cfg.AnalyticExpr{
		Name:    "expr1",
		Expr:    "(M1 + M2) / 2",
		Metrics: []string{"M1", "M2"},
	}
	expr2 := cfg.AnalyticExpr{
		Name:        "expr2",
		Expr:        "M1 / C1 + M2 / C2",
		Metrics:     []string{"M1", "M2"},
		Constraints: []string{"C1", "C2"},
	}
	expr3 := cfg.AnalyticExpr{
		Name:        "expr3",
		Expr:        "M1 * C1 + M3 / C3",
		Metrics:     []string{"M1", "M3"},
		Constraints: []string{"C1", "C3"},
	}

	expressions := map[string]cfg.AnalyticExpr{
		"expr1": expr1,
		"expr2": expr2,
		"expr3": expr3,
	}

	cfg.SetAnalyticExpr(expressions)
}
