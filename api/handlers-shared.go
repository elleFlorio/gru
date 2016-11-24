package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/elleFlorio/gru/data"
)

// /gru/v1/shared/
func GetSharedData(w http.ResponseWriter, r *http.Request) {
	info, _ := data.GetSharedLocal()

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

func PostSharedData(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	shared, err := readShared(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data.AddFriendData(r.RemoteAddr, shared)

	w.WriteHeader(http.StatusAccepted)
}

func readShared(r *http.Request) (data.Shared, error) {
	var err error
	var shared data.Shared

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.WithField("err", err).Errorln("Error reading command body")
		return data.Shared{}, err
	}

	if err = r.Body.Close(); err != nil {
		log.WithField("err", err).Errorln("Error closing command body")
		return data.Shared{}, err
	}

	if err = json.Unmarshal(body, &shared); err != nil {
		log.WithField("err", err).Errorln("Error unmarshaling command body")
		return data.Shared{}, err
	}

	return shared, nil
}
