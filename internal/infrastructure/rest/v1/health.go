package v1

import (
	"log"
	"net/http"
)

func (restApiV1 *RestApiV1) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("OK"))
	if err != nil {
		log.Fatal("could not write health response")
	}
}
