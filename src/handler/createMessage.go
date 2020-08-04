package handler

import (
	"net/http"
	"request-golang/src/service"
	"request-golang/src/util"
)

// CreateMessage is a handler for route GET /api/message
func CreateMessage(datastore service.MessageDatastore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var i service.Message

		// Validate the JSON.
		statusCode, err := util.Unmarshal(w, r, &i)
		if err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}

		// Check if the email is a valid email.
		if err = validateEmail(i.Email); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Create the message in the database.
		if err = datastore.CreateMessage(i); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
