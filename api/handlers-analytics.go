package api

import (
	"encoding/json"
	"net/http"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/gorilla/mux"

	"github.com/elleFlorio/gru/autonomic/analyzer"
)

// /gru/v1/analytics/
func GetAnalyticsNode(w http.ResponseWriter, r *http.Request) {
	analytics := analyzer.GetAnalytics()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(analytics); err != nil {
		log.WithFields(log.Fields{
			"status":  "http response",
			"request": "GetAnalyticsNode",
			"error":   err,
		}).Errorln("API Server")
	}
}

// /gru/v1/analytics/services
func GetAnalyticsServices(w http.ResponseWriter, r *http.Request) {
	analytics := analyzer.GetServicesAnalytics()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(analytics); err != nil {
		log.WithFields(log.Fields{
			"status":  "http response",
			"request": "GetAnalyticsServices",
			"error":   err,
		}).Errorln("API Server")
	}
}

// /gru/v1/analytics/services/{name}
func GetAnalyticsService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	analytics := analyzer.GetServiceAnalytics(name)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(analytics); err != nil {
		log.WithFields(log.Fields{
			"status":  "http response",
			"request": "GetAnalyticsService",
			"error":   err,
		}).Errorln("API Server")
	}
}

// /gru/v1/analytics/system
func GetAnalyticsSystem(w http.ResponseWriter, r *http.Request) {
	analytics := analyzer.GetSystemAnalytics()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(analytics); err != nil {
		log.WithFields(log.Fields{
			"status":  "http response",
			"request": "GetAnalyticsSystem",
			"error":   err,
		}).Errorln("API Server")
	}
}
