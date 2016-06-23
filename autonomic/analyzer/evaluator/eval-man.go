package evaluator

import (
	"strconv"
	"strings"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/soniah/evaler"

	cfg "github.com/elleFlorio/gru/configuration"
	srv "github.com/elleFlorio/gru/service"
)

func ComputeMetricAnalytics(service string, metrics map[string]float64) map[string]float64 {
	expressions := cfg.GetExpr()
	srvExprList := srv.GetServiceExpressionsList(service)
	metricAnalytics := make(map[string]float64, len(expressions))

	for _, expr := range srvExprList {
		if curExpr, ok := expressions[expr]; ok {
			toEval := buildExpression(curExpr.Expr, metrics)
			result, err := evaler.Eval(toEval)
			if err != nil {
				log.WithFields(log.Fields{
					"err":  err,
					"expr": toEval,
				}).Errorln("Error evaluating expression")

				metricAnalytics[expr] = 0.0
			} else {
				metricAnalytics[expr] = evaler.BigratToFloat(result)
			}
		} else {
			log.WithField("expr", expr).Errorln("Cannot compute expression: expression unknown")
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
