package handler

import (
	"net/http"
	"request-golang/src/config"
	"request-golang/src/service"
	"request-golang/src/util"
)

type SendMessagesInput struct {
	MagicNumber int `json:"magic_number"`
}

func SendMessages(datastore service.MessageDatastore, c *config.SmtpConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate the JSON and get the data.
		var i SendMessagesInput
		errStatusCode, err := util.Unmarshal(w, r, &i)
		if err != nil {
			http.Error(w, err.Error(), errStatusCode)
			return
		}

		// Send messages with given magic number and delete them from the database
		if err = datastore.SendMessages(i.MagicNumber, c); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
