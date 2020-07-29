package handler

import (
	"github.com/gorilla/mux"
	"net/http"
)

func GetMessagesByEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		email := vars["email"]

	}
}

