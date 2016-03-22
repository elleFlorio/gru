package planner

import (
	"encoding/json"
	"fmt"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/autonomic/planner/policy"
	"github.com/elleFlorio/gru/autonomic/planner/strategy"
	"github.com/elleFlorio/gru/data"
	"github.com/elleFlorio/gru/enum"
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

func Run(analytics data.GruAnalytics) *policy.Policy {
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

	}

	return chosenPolicy
}

// TODO I should remove this and start from the stats to compute something
// for the services that are not created yet. In that way I just need to call
// service.List()
func getServicesListFromAnalytics(analytics data.GruAnalytics) []string {
	list := make([]string, 0, len(analytics.Service))
	for srv, _ := range analytics.Service {
		list = append(list, srv)
	}

	return list
}

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
	targets := make([]string, 0, len(chosenPolicy.Targets))
	for _, target := range chosenPolicy.Targets {
		targets = append(targets, target)
	}

	log.WithFields(log.Fields{
		"name":    chosenPolicy.Name,
		"weight":  fmt.Sprintf("%.2f", chosenPolicy.Weight),
		"targets": targets,
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
