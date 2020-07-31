package handler

import (
	"net/http"
	"request-golang/service"
	"request-golang/util"
)

func CreateMessage(datastore service.MessageDatastore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//bodyBytes, err := ioutil.ReadAll(r.Body)
		//defer r.Body.Close()
		//if err != nil {
		//	w.WriteHeader(http.StatusInternalServerError)
		//	w.Write([]byte(err.Error()))
		//}

		var i service.Message
		statusCode, err := util.Unmarshal(w, r, i)
		if err != nil {
			w.WriteHeader(statusCode)
			w.Write([]byte(err.Error()))
		}
		//if err = json.Unmarshal(bodyBytes, &i); err != nil {
		//	w.WriteHeader(http.StatusBadRequest)
		//	w.Write([]byte(err.Error()))
		//}
		if err = datastore.CreateMessage(i); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		w.WriteHeader(http.StatusCreated)
	}
}
