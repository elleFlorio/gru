package api

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	"github.com/elleFlorio/gru/autonomic/analyzer"
)

// /gru/v1/analytics/services
func GetAnalyticsServices(w http.ResponseWriter, r *http.Request) {
	analytics := analyzer.GetServicesAanalytics()

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
	analytics := analyzer.GetServiceAanalytics(name)

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

// /gru/v1/analytics/instances
func GetAnalyticsInstances(w http.ResponseWriter, r *http.Request) {
	analytics := analyzer.GetInstancesAanalytics()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(analytics); err != nil {
		log.WithFields(log.Fields{
			"status":  "http response",
			"request": "GetAnalyticsInstances",
			"error":   err,
		}).Errorln("API Server")
	}
}

// /gru/v1/analytics/instances/{id}
func GetAnalyticsInstance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	analytics := analyzer.GetInstanceAanalytics(id)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(analytics); err != nil {
		log.WithFields(log.Fields{
			"status":  "http response",
			"request": "GetAnalyticsInstance",
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
