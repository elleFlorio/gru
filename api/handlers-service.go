package api

import (
	"encoding/json"
	"net/http"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	cfg "github.com/elleFlorio/gru/configuration"
)

// /gru/v1/services
func GetInfoServices(w http.ResponseWriter, r *http.Request) {
	services := cfg.GetServices()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(services); err != nil {
		log.WithFields(log.Fields{
			"status":  "http response",
			"request": "GetInfoServices",
			"error":   err,
		}).Errorln("API Server")
	}
}
