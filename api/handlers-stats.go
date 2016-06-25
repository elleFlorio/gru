package api

import (
	"encoding/json"
	"net/http"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/gorilla/mux"

	"github.com/elleFlorio/gru/autonomic/monitor"
	mtr "github.com/elleFlorio/gru/autonomic/monitor/metric"
)

// /gru/v1/stats
func GetStatsNode(w http.ResponseWriter, r *http.Request) {
	stats := monitor.GetStats()

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

// /gru/v1/stats/user/{service}/{metric}
func PostServiceMetrics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service := vars["service"]
	metric := vars["metric"]
	var userValue struct {
		values []float64
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userValue); err != nil {
		log.WithFields(log.Fields{
			"status":  "http post",
			"request": "PostUpdateMetrics",
			"error":   err,
		}).Errorln("API Server")
		w.WriteHeader(http.StatusBadRequest)
	} else {
		mtr.UpdateUserMetric(service, metric, userValue.values)
		w.WriteHeader(http.StatusOK)
	}
}
