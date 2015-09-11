package api

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	"github.com/elleFlorio/gru/autonomic/monitor"
)

// /gru/v1/stats
func GetStatsNode(w http.ResponseWriter, r *http.Request) {
	stats := monitor.GetNodeStats()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		log.WithFields(log.Fields{
			"status":  "http response",
			"request": "GetStatsNode",
			"error":   err,
		}).Errorln("API Server")
	}
}

// /gru/v1/stats/services
func GetStatsServices(w http.ResponseWriter, r *http.Request) {
	stats := monitor.GetServicesStats()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		log.WithFields(log.Fields{
			"status":  "http response",
			"request": "GetStatsServices",
			"error":   err,
		}).Errorln("API Server")
	}
}

// /gru/v1/stats/services/{name}
func GetStatsService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	stats := monitor.GetServiceStats(name)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		log.WithFields(log.Fields{
			"status":  "http response",
			"request": "GetStatsService",
			"error":   err,
		}).Errorln("API Server")
	}
}

// /gru/v1/stats/instances
func GetStatsInstances(w http.ResponseWriter, r *http.Request) {
	stats := monitor.GetInstancesStats()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		log.WithFields(log.Fields{
			"status":  "http response",
			"request": "GetStatsInstances",
			"error":   err,
		}).Errorln("API Server")
	}
}

// /gru/v1/stats/instances/{id}
func GetStatsInstance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	stats := monitor.GetInstanceStats(id)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		log.WithFields(log.Fields{
			"status":  "http response",
			"request": "GetStatsInstance",
			"error":   err,
		}).Errorln("API Server")
	}
}

// /gru/v1/stats/system
func GetStatsSystem(w http.ResponseWriter, r *http.Request) {
	stats := monitor.GetSystemStats()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		log.WithFields(log.Fields{
			"status":  "http response",
			"request": "GetStatsSystem",
			"error":   err,
		}).Errorln("API Server")
	}
}
