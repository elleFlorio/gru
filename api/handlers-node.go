package api

import (
	"encoding/json"
	"net/http"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/node"
)

// /gru/v1/node
func GetInfoNode(w http.ResponseWriter, r *http.Request) {
	info := node.GetNode()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(info); err != nil {
		log.WithFields(log.Fields{
			"status":  "http response",
			"request": "GetInfoNode",
			"error":   err,
		}).Errorln("API Server")
	}
}
