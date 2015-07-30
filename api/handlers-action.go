package api

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/action"
)

// /gru/v1/actions
func GetInfoActions(w http.ResponseWriter, r *http.Request) {
	actions := action.List()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(actions); err != nil {
		log.WithFields(log.Fields{
			"status":  "http response",
			"request": "GetInfoPolicies",
			"error":   err,
		}).Errorln("API Server")
	}
}
