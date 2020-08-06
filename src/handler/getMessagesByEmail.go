package handler

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"net/http"
	"request-golang/src/service"
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

		// Try to transform the string limit to an integer and check if it's not negative.
		// If its value is 0, we won't apply the limit, since it's also a default "undefined" value.
		limit, err := strconv.Atoi(strLimit)
		if (err != nil && strLimit != "") || limit < 0 {
			http.Error(w, "limit parameter invalid: it must be a positive integer", http.StatusBadRequest)
			return
		}

		// It will get validated later.
		cursor := r.URL.Query().Get(CURSOR)

		// Define the input data for the service.
		i := service.GetMessagesByEmailInput{
			Email:  email,
			Params: service.FetchParams{
				Limit: limit,
				Cursor: cursor,
			},
		}

		// Validate the input. The validation rules are specified in the struct definition.
		if err := validator.New().Struct(i); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		messages, endCursor := datastore.GetMessagesByEmail(i)
		resData := getMessagesByEmailResponse{
			Messages:  messages,
			EndCursor: endCursor,
		}

		payload, err := json.Marshal(resData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(payload)
	}
}
