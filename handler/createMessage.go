package handler

import (
	"net/http"
	"request-golang/service"
)

func CreateMessage(datastore service.MessageDatastore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

