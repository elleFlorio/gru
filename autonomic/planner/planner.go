package planner

import (
	"fmt"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/autonomic/planner/policy"
	"github.com/elleFlorio/gru/autonomic/planner/strategy"
	"github.com/elleFlorio/gru/data"
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

func Run(clusterData data.Shared) *data.Policy {
	log.WithField("status", "init").Debugln("Gru Planner")
	defer log.WithField("status", "done").Debugln("Gru Planner")

	var chosenPolicy *data.Policy

	if len(clusterData.Service) == 0 {
		log.Warnln("No cluster data for policy computation")
		return chosenPolicy
	}

	srvList := getServicesListFromClusterData(clusterData)
	policies := policy.CreatePolicies(srvList, clusterData)
	chosenPolicy = currentStrategy.MakeDecision(policies)
	data.SavePolicy(*chosenPolicy)
	displayPolicy(chosenPolicy)

	return chosenPolicy
}

func getServicesListFromClusterData(clusterData data.Shared) []string {
	list := make([]string, 0, len(clusterData.Service))
	for srv, _ := range clusterData.Service {
		list = append(list, srv)
	}

	return list
}

func displayPolicy(chosenPolicy *data.Policy) {
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
