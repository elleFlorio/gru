package planner

import (
	"encoding/json"

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

	analytics, err := analyzer.GetAnalyzerData()
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

func buildPlans(analytics analyzer.GruAnalytics) []strategy.GruPlan {
	plans := []strategy.GruPlan{}
	plans = append(plans, createNoActionPlan())

	if len(analytics.Service) == 0 {
		log.Warnln("No service for building plans.")
		return plans
	}

	policies := policy.GetPolicies()
	for _, name := range service.List() {
		for _, plc := range policies {
			label := plc.Label(name, analytics)
			target, _ := service.GetServiceByName(name)
			actions := plc.Actions()
			plan := strategy.GruPlan{label, target, actions}
			log.WithFields(log.Fields{
				"policy":  plc.Name(),
				"label":   label.ToString(),
				"service": name,
			}).Debugln("Computed policy")
			plans = append(plans, plan)
		}
	}

	return plans
}

func createNoActionPlan() strategy.GruPlan {
	srv := service.Service{Name: "NoAction"}
	return strategy.GruPlan{enum.GREEN, &srv, []enum.Action{enum.NOACTION}}
}

func savePlan(plan *strategy.GruPlan) error {
	data, err := convertPlanToData(*plan)
	if err != nil {
		log.WithField("error", "Cannot convert plan to data").Debugln("Running Planner")
		return err
	} else {
		storage.StoreLocalData(data, enum.PLANS)
	}

	return nil
}

func convertPlanToData(plan strategy.GruPlan) ([]byte, error) {
	data, err := json.Marshal(plan)
	if err != nil {
		log.WithField("error", err).Errorln("Error marshaling plan data")
		return nil, err
	}

	return data, nil
}

func GetPlannerData() (strategy.GruPlan, error) {
	plan := strategy.GruPlan{}
	dataPlan, err := storage.GetLocalData(enum.PLANS)
	if err != nil {
		log.WithField("error", err).Errorln("Cannot retrieve plan data.")
	} else {
		plan, err = convertDataToPlan(dataPlan)
	}

	return plan, err
}

func convertDataToPlan(data []byte) (strategy.GruPlan, error) {
	plan := strategy.GruPlan{}
	err := json.Unmarshal(data, &plan)
	if err != nil {
		log.WithField("error", err).Errorln("Error unmarshaling plan data")
	}

	return plan, err
}
