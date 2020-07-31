package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"request-golang/service"
)

type SendMessageInput struct {
	MagicNumber int `json:"magic_number"`
}

func SendMessage(datastore service.MessageDatastore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		var i SendMessageInput
		if err = json.Unmarshal(bodyBytes, &i); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		}

		if err = datastore.SendMessages(i.MagicNumber); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		w.WriteHeader(http.StatusOK)
	}
}
