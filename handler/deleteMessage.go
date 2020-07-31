package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"request-golang/service"
)

type DeleteMessageInput struct {
	MagicNumber int `json:"magic_number"`
}

func DeleteMessage(datastore service.MessageDatastore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		var i DeleteMessageInput
		if err = json.Unmarshal(bodyBytes, &i); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		}

		if err = datastore.DeleteMessage(i.MagicNumber); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		w.WriteHeader(http.StatusOK)
	}
}

