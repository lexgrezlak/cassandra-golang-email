package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"request-golang/service"
)

func GetMessagesByEmail(datastore service.MessageDatastore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		email := vars["email"]
		messages, err := datastore.GetMessagesByEmail(email)
		if err != nil {
			fmt.Errorf("failed to get all messages by email: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		if payload, err := json.Marshal(messages); err == nil {
			w.Write(payload)
		}
	}
}

