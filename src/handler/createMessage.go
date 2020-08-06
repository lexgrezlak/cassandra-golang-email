package handler

import (
	"github.com/go-playground/validator/v10"
	"net/http"
	"request-golang/src/service"
	"request-golang/src/util"
)

// CreateMessage is a handler for route GET /api/message
func CreateMessage(datastore service.MessageDatastore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var i service.CreateMessageInput

		// Validate and unmarshal the JSON.
		statusCode, err := util.Unmarshal(w, r, &i)
		if err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}

		// Validate the input.
		if err = validator.New().Struct(i); err != nil {
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
