package planner

import (
	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/autonomic/analyzer"
	"github.com/elleFlorio/gru/autonomic/planner/policy"
	"github.com/elleFlorio/gru/autonomic/planner/strategy"
	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/service"
	"github.com/elleFlorio/gru/storage"
)

var currentStrategy strategy.GruStrategy

func SetPlannerStrategy(strategyName string) {
	strtg, err := strategy.New(strategyName)
	if err != nil {
		log.WithField("error", err).Errorln("Strategy cannot be set")

		// If error use default one
		strtg, err = strategy.New("dummy")
	}

	currentStrategy = strtg

	log.WithField("strategy", strtg.Name()).Infoln("Strategy initialized")
}

func Run() {
	log.WithField("status", "start").Debugln("Running planner")
	defer log.WithField("status", "done").Debugln("Running planner")

	analytics, err := retrieveAnalytics()
	if err != nil {
		log.WithField("error", "Cannot compute plans").Errorln("Running Planner.")
	} else {
		plans := buildPlans(analytics)
		thePlan := currentStrategy.MakeDecision(plans)
		err := savePlan(thePlan)
		if err != nil {
			log.WithField("error", "Plan data not saved ").Errorln("Running Planner")
		}
	}
}

func retrieveAnalytics() (analyzer.GruAnalytics, error) {
	analytics := analyzer.GruAnalytics{}
	dataAnalyics, err := storage.GetData(enum.CLUSTER.ToString(), enum.ANALYTICS)
	if err != nil {
		log.WithField("error", err).Errorln("Cannot retrieve analytics data.")
	} else {
		analytics, err = analyzer.ConvertDataToAnalytics(dataAnalyics)
	}

	return analytics, err
}

func buildPlans(analytics analyzer.GruAnalytics) []strategy.GruPlan {
	plans := []strategy.GruPlan{}
	policies := policy.GetPolicies()

	for _, name := range service.List() {
		for _, plc := range policies {
			label := plc.Label(name, analytics)
			target, _ := service.GetServiceByName(name)
			actions := plc.Actions()
			plan := strategy.GruPlan{label, target, actions}

			plans = append(plans, plan)
		}
	}

	return plans
}

func savePlan(plan *strategy.GruPlan) error {
	data, err := strategy.ConvertPlanToData(*plan)
	if err != nil {
		log.WithField("error", "Cannot convert plan to data").Debugln("Running Planner")
		return err
	} else {
		storage.StoreLocalData(data, enum.PLANS)
	}

	return nil
}
