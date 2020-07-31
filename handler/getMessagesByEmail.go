package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"request-golang/service"
	"strconv"
)

const (
	EMAIL  = "email"
	LIMIT  = "limit"
	CURSOR = "cursor"
)

type getMessagesByEmailResponse struct {
	Messages  []*service.Message `json:"messages"`
	EndCursor string             `json:"endCursor"`
}

func GetMessagesByEmail(datastore service.MessageDatastore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		email := vars[EMAIL]
		strLimit := r.URL.Query().Get(LIMIT)
		limit, err := strconv.Atoi(strLimit)
		if err != nil && strLimit != "" && limit > 1 {
			http.Error(w, "limit parameter must be a positive integer greater than 1", http.StatusBadRequest)
			return
		}
		cursor := r.URL.Query().Get(CURSOR)
		messages, endCursor, err := datastore.GetMessagesByEmail(email, limit, cursor)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		resData := getMessagesByEmailResponse{
			Messages:  messages,
			EndCursor: endCursor,
		}

		payload, err := json.Marshal(resData)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		}
		w.WriteHeader(http.StatusOK)
		w.Write(payload)
	}
}
