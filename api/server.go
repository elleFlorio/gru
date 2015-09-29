package api

import (
	"net/http"

	log "github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

func StartServer(port string) {

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":"+port, router))
}
