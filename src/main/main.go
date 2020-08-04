package main

import (
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
	// Initialize config. We could also set the config path as the environment variable.
	// Environment variables will overwrite the config.
	c, err := config.GetConfig("config.yml")
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
		log.Fatalf("failed to create session: %v", err)
	}
	defer session.Close()

	api := service.NewAPI(session)

	// Set up handlers.
	r := mux.NewRouter()
	r.HandleFunc("/api/message", handler.CreateMessage(api)).Methods("POST")
	r.HandleFunc("/api/send", handler.SendMessages(api, c.Smtp)).Methods("POST")
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
	log.Printf("Listening at: %v", srv.Addr)
	err = srv.ListenAndServe()
	log.Fatal(err)
}
