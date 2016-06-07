package evaluator

import (
	"strconv"
	"strings"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/soniah/evaler"

	cfg "github.com/elleFlorio/gru/configuration"
)

func ComputeMetricAnalytics(metrics map[string]float64) map[string]float64 {
	expressions := cfg.GetExpr()
	metricAnalytics := make(map[string]float64, len(expressions))

	for exprName, expr := range expressions {
		toEval := buildExpression(expr.Expr, metrics)
		result, err := evaler.Eval(toEval)
		if err != nil {
			log.WithFields(log.Fields{
				"err":  err,
				"expr": toEval,
			}).Errorln("Error evaluating expression")

			metricAnalytics[exprName] = 0.0
		} else {
			metricAnalytics[exprName] = evaler.BigratToFloat(result)
		}
	}

	return metricAnalytics
}

func buildExpression(expr string, metrics map[string]float64) string {
	toBuild := expr
	for varName, varValue := range metrics {
		toBuild = strings.Replace(toBuild, varName, strconv.FormatFloat(varValue, 'f', -1, 64), -1)
	}

	return toBuild
}
