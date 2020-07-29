package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"request-golang/service"
)

func CreateMessage(datastore service.MessageDatastore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var i service.Message
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err) // TODO
		}
		if err = json.Unmarshal(body, &i); err != nil {
			panic(err) // TODO
		}
		if err = datastore.CreateMessage(i); err != nil {
			fmt.Errorf("failed to create message: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

