package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"request-golang/handler"
	"time"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/message", handler.CreateMessage()).Methods("POST")
	r.HandleFunc("/api/send", handler.DeleteMessage()).Methods("POST")
	r.HandleFunc("/api/messages/{email}", handler.GetMessagesByEmail()).Methods("GET")

	srv := &http.Server{
		Handler: r,
		Addr: "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
	}
	fmt.Println("Listening at:", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}