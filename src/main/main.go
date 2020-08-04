package main

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"request-golang/src/config"
	"request-golang/src/handler"
	"request-golang/src/service"
	"time"
)



func main() {
	// Initialize config.
	c, err := config.LoadConfig("../config.yml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Set up the cluster and create a session.
	cluster := gocql.NewCluster(c.Db.Host)
	cluster.Keyspace = c.Db.Keyspace
	cluster.Consistency = gocql.All
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: c.Db.Username,
		Password: c.Db.Password,
	}
	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	api := service.NewAPI(session)

	// Set up handlers.
	r := mux.NewRouter()
	r.HandleFunc("/api/message", handler.CreateMessage(api)).Methods("POST")
	r.HandleFunc("/api/send", handler.SendMessages(api)).Methods("POST")
	// For paginated results use ?limit=5&cursor=hello-world, for example.
	r.HandleFunc("/api/messages/{email}", handler.GetMessagesByEmail(api)).Methods("GET")

	// Set up the server.
	srv := &http.Server{
		Handler:      r,
		Addr:         c.Server.Address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	fmt.Println("Listening at:", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
