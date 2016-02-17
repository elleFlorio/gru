package planner

import (
	"encoding/json"
	"fmt"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/autonomic/analyzer"
	"github.com/elleFlorio/gru/autonomic/planner/policy"
	"github.com/elleFlorio/gru/autonomic/planner/strategy"
	// cfg "github.com/elleFlorio/gru/configuration"
	"github.com/elleFlorio/gru/enum"
	// "github.com/elleFlorio/gru/service"
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

func Run(analytics analyzer.GruAnalytics) *policy.Policy {
	log.WithField("status", "init").Debugln("Gru Planner")
	defer log.WithField("status", "done").Debugln("Gru Planner")
	var chosenPolicy *policy.Policy

	if len(analytics.Service) == 0 {
		log.WithField("err", "No services analytics").Warnln("Cannot compute plans.")
	} else {
		srvList := getServicesListFromAnalytics(analytics)
		policies := policy.CreatePolicies(srvList, analytics)
		chosenPolicy = currentStrategy.MakeDecision(policies)
		err := savePolicy(chosenPolicy)
		if err != nil {
			log.WithField("err", err).Errorln("Planner data not saved")
		}

		displayPolicy(chosenPolicy)

		// plans := buildPlans(analytics)
		// thePlan = currentStrategy.MakeDecision(plans)
		// err := savePlan(thePlan)
		// if err != nil {
		// 	log.WithField("err", err).Errorln("Plan data not saved")
		// }

		// displayThePlan(thePlan)
	}

	return chosenPolicy
}

// TODO I should remove this and start from the stats to compute something
// for the services that are not created yet. In that way I just need to call
// service.List()
func getServicesListFromAnalytics(analytics analyzer.GruAnalytics) []string {
	list := make([]string, 0, len(analytics.Service))
	for srv, _ := range analytics.Service {
		list = append(list, srv)
	}

	return list
}

// func buildPlans(analytics analyzer.GruAnalytics) []strategy.GruPlan {
// 	plans := []strategy.GruPlan{}

// 	if len(analytics.Service) == 0 {
// 		log.Warnln("No service for building plans.")
// 		noServicePlan := strategy.GruPlan{
// 			"noaction",
// 			1.0,
// 			&cfg.Service{Name: "noService"},
// 			[]enum.Action{enum.NOACTION},
// 		}
// 		plans = append(plans, noServicePlan)
// 		return plans
// 	}

// 	policies := policy.GetPolicies()
// 	weight_max := 0.0
// 	for _, name := range service.List() {
// 		for _, plc := range policies {
// 			weight := plc.Weight(name, analytics)
// 			target, _ := service.GetServiceByName(name)
// 			actions := plc.Actions()
// 			plan := strategy.GruPlan{plc.Name(), weight, target, actions}
// 			log.WithFields(log.Fields{
// 				"policy":  plc.Name(),
// 				"weight":  weight,
// 				"service": name,
// 			}).Debugln("Computed policy")
// 			plans = append(plans, plan)

// 			if weight >= weight_max {
// 				weight_max = weight
// 			}
// 		}
// 	}
// 	weight_na := 1 - weight_max
// 	plan_na := strategy.GruPlan{
// 		"noaction",
// 		weight_na,
// 		&cfg.Service{Name: "NoService"},
// 		[]enum.Action{enum.NOACTION},
// 	}
// 	log.WithFields(log.Fields{
// 		"policy":  "NoAction",
// 		"weight":  weight_na,
// 		"service": "NoService",
// 	}).Debugln("Computed policy")
// 	plans = append(plans, plan_na)

// 	return plans
// }

func savePolicy(chosenPolicy *policy.Policy) error {
	data, err := convertPolicyToData(chosenPolicy)
	if err != nil {
		log.WithField("err", err).Debugln("Cannot convert policy to data")
		return err
	} else {
		storage.StoreLocalData(data, enum.POLICIES)
	}

	return nil
}

func convertPolicyToData(chosenPolicy *policy.Policy) ([]byte, error) {
	data, err := json.Marshal(*chosenPolicy)
	if err != nil {
		log.WithField("error", err).Errorln("Error marshaling plan data")
		return nil, err
	}

	return data, nil
}

func displayPolicy(chosenPolicy *policy.Policy) {
	log.WithFields(log.Fields{
		"name":    chosenPolicy.Name,
		"weight":  fmt.Sprintf("%.2f", chosenPolicy.Weight),
		"targets": chosenPolicy.Targets,
	}).Infoln("Policy to actuate")
}

func GetPlannerData() (policy.Policy, error) {
	plc := policy.Policy{}
	dataPlan, err := storage.GetLocalData(enum.POLICIES)
	if err != nil {
		log.WithField("err", err).Warnln("Cannot retrieve plan data")
	} else {
		plc, err = convertDataToPolicy(dataPlan)
	}

	return plc, err
}

func convertDataToPolicy(data []byte) (policy.Policy, error) {
	plc := policy.Policy{}
	err := json.Unmarshal(data, &plc)
	if err != nil {
		log.WithField("err", err).Errorln("Error unmarshaling plan data")
		return policy.Policy{}, err
	}

	return plc, nil
}
