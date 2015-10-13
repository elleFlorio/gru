package strategy

import (
	"encoding/json"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/service"
)

type GruPlan struct {
	Label   enum.Label
	Target  *service.Service
	Actions []enum.Action
}

func ConvertPlanToData(plan GruPlan) ([]byte, error) {
	data, err := json.Marshal(plan)
	if err != nil {
		log.WithField("error", err).Errorln("Error marshaling analytics data")
		return nil, err
	}

	return data, nil
}

func ConvertDataToPlan(data []byte) (GruPlan, error) {
	plan := GruPlan{}
	err := json.Unmarshal(data, &plan)
	if err != nil {
		log.WithField("error", err).Errorln("Error unmarshaling analytics data")
	}

	return plan, err
}
