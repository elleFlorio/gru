package planner

import (
	"encoding/json"
	"fmt"

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
		log.WithField("err", err).Errorln("Strategy cannot be set")

		// If error use default one
		strtg, err = strategy.New("dummy")
	}

	currentStrategy = strtg

	log.WithField("strategy", strtg.Name()).Infoln("Strategy initialized")
}

func Run(analytics analyzer.GruAnalytics) *strategy.GruPlan {
	log.WithField("status", "init").Debugln("Gru Planner")
	defer log.WithField("status", "done").Debugln("Gru Planner")
	var thePlan *strategy.GruPlan

	if len(analytics.Service) == 0 {
		log.WithField("err", "No services analytics").Errorln("Cannot compute plans.")
	} else {
		plans := buildPlans(analytics)
		thePlan = currentStrategy.MakeDecision(plans)
		err := savePlan(thePlan)
		if err != nil {
			log.WithField("err", err).Errorln("Plan data not saved")
		}

		displayThePlan(thePlan)
	}

	return thePlan
}

func buildPlans(analytics analyzer.GruAnalytics) []strategy.GruPlan {
	plans := []strategy.GruPlan{}

	if len(analytics.Service) == 0 {
		log.Warnln("No service for building plans.")
		noServicePlan := strategy.GruPlan{
			"noaction",
			1.0,
			&service.Service{Name: "noService"},
			[]enum.Action{enum.NOACTION},
		}
		plans = append(plans, noServicePlan)
		return plans
	}

	policies := policy.GetPolicies()
	weight_max := 0.0
	for _, name := range service.List() {
		for _, plc := range policies {
			weight := plc.Weight(name, analytics)
			target, _ := service.GetServiceByName(name)
			actions := plc.Actions()
			plan := strategy.GruPlan{plc.Name(), weight, target, actions}
			log.WithFields(log.Fields{
				"policy":  plc.Name(),
				"weight":  weight,
				"service": name,
			}).Debugln("Computed policy")
			plans = append(plans, plan)

			if weight >= weight_max {
				weight_max = weight
			}
		}
	}
	weight_na := 1 - weight_max
	plan_na := strategy.GruPlan{
		"noaction",
		weight_na,
		&service.Service{Name: "NoService"},
		[]enum.Action{enum.NOACTION},
	}
	log.WithFields(log.Fields{
		"policy":  "NoAction",
		"weight":  weight_na,
		"service": "NoService",
	}).Debugln("Computed policy")
	plans = append(plans, plan_na)

	return plans
}

func savePlan(plan *strategy.GruPlan) error {
	data, err := convertPlanToData(*plan)
	if err != nil {
		log.WithField("err", err).Debugln("Cannot convert plan to data")
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

func displayThePlan(thePlan *strategy.GruPlan) {
	log.WithFields(log.Fields{
		"target":  thePlan.Target.Name,
		"actions": thePlan.Actions.ToString(),
		"weight":  fmt.Sprintf("%.2f", thePlan.Weight),
	}).Infoln("Plan computed")
}

func GetPlannerData() (strategy.GruPlan, error) {
	plan := strategy.GruPlan{}
	dataPlan, err := storage.GetLocalData(enum.PLANS)
	if err != nil {
		log.WithField("err", err).Errorln("Cannot retrieve plan data.")
	} else {
		plan, err = convertDataToPlan(dataPlan)
	}

	return plan, err
}

func convertDataToPlan(data []byte) (strategy.GruPlan, error) {
	plan := strategy.GruPlan{}
	err := json.Unmarshal(data, &plan)
	if err != nil {
		log.WithField("err", err).Errorln("Error unmarshaling plan data")
	}

	return plan, err
}
