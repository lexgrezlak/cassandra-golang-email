package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"request-golang/service"
	"strconv"
)

const (
	EMAIL = "email"
	LIMIT = "limit"
	CURSOR = "cursor"
)

type getMessagesByEmailResponse struct {
	Messages []*service.Message `json:"messages"`
	EndCursor string `json:"endCursor"`
}

func GetMessagesByEmail(datastore service.MessageDatastore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		email := vars[EMAIL]
		strLimit := r.URL.Query().Get(LIMIT)
		limit, err := strconv.Atoi(strLimit)
		if err != nil && strLimit != "" {
			http.Error(w, "limit parameter is not a number", http.StatusBadRequest)
			return
		}
		cursor := r.URL.Query().Get(CURSOR)
		messages, endCursor, err := datastore.GetMessagesByEmail(email, limit, cursor)
		if err != nil {
			fmt.Errorf("failed to get all messages by email: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		fmt.Printf("messages %v, cursor %v", messages, endCursor)
		resData := getMessagesByEmailResponse{
			Messages:  messages,
			EndCursor: endCursor,
		}
		w.WriteHeader(http.StatusOK)
		if payload, err := json.Marshal(resData); err == nil {
			w.Write(payload)
		}
	}
}

