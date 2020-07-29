package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"request-golang/service"
)

func GetMessagesByEmail(datastore service.MessageDatastore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		email := vars["email"]
		fmt.Printf("email: %v", email)

	}
}

