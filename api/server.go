package api

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func StartServer(port string) {

	router := NewRouter()

	log.Fatal(http.ListenAndServe(port, router))
}
