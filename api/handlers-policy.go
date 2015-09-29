package api

import (
	"encoding/json"
	"net/http"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/gorilla/mux"

	"github.com/elleFlorio/gru/policy"
)

type plc struct {
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Level        string   `json:"level"`
	Target       string   `json:"target"`
	TargetStatus string   `json:"targetstatus"`
	Actions      []string `json:"actions"`
}

// /gru/v1/policies
func GetInfoPolicies(w http.ResponseWriter, r *http.Request) {
	policies := policy.GetPolicies("proactive")
	policies = append(policies, policy.GetPolicies("reactive")...)
	plcs := createPoliciesJson(policies)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(plcs); err != nil {
		log.WithFields(log.Fields{
			"status":  "http response",
			"request": "GetInfoPolicies",
			"error":   err,
		}).Errorln("API Server")
	}
}

// /gru/v1/policies/{type}
func GetInfoPoliciesByType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pType := vars["type"]
	policies := policy.GetPolicies(pType)
	plcs := createPoliciesJson(policies)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(plcs); err != nil {
		log.WithFields(log.Fields{
			"status":  "http response",
			"request": "GetInfoPoliciesByType",
			"error":   err,
		}).Errorln("API Server")
	}
}

func createPoliciesJson(policies []policy.GruPolicy) []plc {
	plcs := make([]plc, 0, len(policies))

	for _, p := range policies {
		plc_tmp := plc{
			p.Name(),
			p.Type(),
			p.Level(),
			p.Target(),
			p.TargetStatus(),
			p.Actions(),
		}
		plcs = append(plcs, plc_tmp)
	}

	return plcs
}
