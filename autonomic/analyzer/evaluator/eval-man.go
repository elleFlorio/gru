package evaluator

import (
	"math"
	"strconv"
	"strings"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/soniah/evaler"

	cfg "github.com/elleFlorio/gru/configuration"
	srv "github.com/elleFlorio/gru/service"
)

func ComputeMetricAnalytics(service string, metrics map[string]float64) map[string]float64 {
	expressions := cfg.GetExpr()
	srvExprList := srv.GetServiceExpressionsList(service)
	srvConstraints := srv.GetServiceConstraints(service)
	metricAnalytics := make(map[string]float64, len(expressions))

	for _, expr := range srvExprList {
		if curExpr, ok := expressions[expr]; ok {
			log.WithField("expr", expr).Debugln("Evaluating expression")
			toEval := buildExpression(curExpr, metrics, srvConstraints)
			result, err := evaler.Eval(toEval)
			if err != nil {
				log.WithFields(log.Fields{
					"err":  err,
					"expr": toEval,
				}).Errorln("Error evaluating expression")

				metricAnalytics[expr] = 0.0
			} else {
				value := evaler.BigratToFloat(result)
				value = math.Min(value, 1.0)
				value = math.Max(value, 0.0)
				metricAnalytics[expr] = evaler.BigratToFloat(result)
			}

			log.WithFields(log.Fields{
				"service": service,
				"expr":    expr,
				"value":   metricAnalytics[expr],
			}).Debugln("Expression evaluated")

		} else {
			log.WithField("expr", expr).Errorln("Cannot compute expression: expression unknown")
		}
	}

	return metricAnalytics
}

func buildExpression(expr cfg.Expression, metrics map[string]float64, constraints map[string]float64) string {
	toBuild := expr.Expr

	for _, metric := range expr.Metrics {
		if value, ok := metrics[metric]; ok {
			toBuild = strings.Replace(toBuild, metric, strconv.FormatFloat(value, 'f', -1, 64), -1)
		} else {
			log.WithField("metric", metric).Errorln("Cannot build expression: metric unknown")
			return "noexp"
		}

	}

	for _, constraint := range expr.Constraints {
		if value, ok := constraints[constraint]; ok {
			toBuild = strings.Replace(toBuild, constraint, strconv.FormatFloat(value, 'f', -1, 64), -1)
		} else {
			log.WithField("constraint", constraint).Errorln("Cannot build expression: constraint unknown")
			return "noexp"
		}

	}

	return toBuild
}
