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
		var i SendMessagesInput
		statusCode, err := util.Unmarshal(w, r, &i)
		if err != nil {
			w.WriteHeader(statusCode)
			w.Write([]byte(err.Error()))
		}

		if err = datastore.SendMessages(i.MagicNumber, c); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		w.WriteHeader(http.StatusOK)
	}
}
