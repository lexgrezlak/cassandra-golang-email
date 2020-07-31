package main

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"request-golang/handler"
	"request-golang/service"
	"time"
)

func main() {
	cluster := gocql.NewCluster("cassandra")
	cluster.Keyspace = "public"
	cluster.Consistency = gocql.All
	cluster.Authenticator = gocql.PasswordAuthenticator{
		// We would normally use environment variables but we're supposed
		// to push the docker image into docker hub and ensure it's working
		// so we're not gonna make others spend time on creating .env file
		// Overall you shouldn't ever store such data in your code,
		// even for testing or development.
		Username: "cassandra",
		Password: "cassandra",
	}
	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	api := service.NewAPI(session)
	r := mux.NewRouter()
	r.HandleFunc("/api/message", handler.CreateMessage(api)).Methods("POST")
	r.HandleFunc("/api/send", handler.SendMessages(api)).Methods("POST")
	// For paginated results use ?limit=5&cursor=hello-world for example
	r.HandleFunc("/api/messages/{email}", handler.GetMessagesByEmail(api)).Methods("GET")

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("Listening at:", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
