package main

import (
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"request-golang/src/config"
	"request-golang/src/handler"
	"request-golang/src/middleware"
	"request-golang/src/service"
	"time"
)

func main() {
	// Initialize config. If it can't find the file, it will load the variables
	// from the environment. It would be a good idea to read the file path to the config
	// from environment, because we might want to have `test.yml` or some other config.
	c, err := config.GetConfig("development.yml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	log.Printf("config has been loaded: %v", c)

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

	// Define the API.
	api := service.NewAPI(session)

	// Define the router.
	r := mux.NewRouter()

	// Set up middleware.
	r.Use(middleware.Logger)
	r.Use(middleware.RequestLimiter)

	// Set up handlers.
	r.HandleFunc("/api/message", handler.CreateMessage(api)).Methods(http.MethodPost)
	r.HandleFunc("/api/send", handler.SendMessages(api, &c.Smtp)).Methods(http.MethodPost)
	// For paginated results use ?limit=5&cursor=hello-world, for example.
	// Limit can be an integer from 1 to 100.
	r.HandleFunc("/api/messages/{email}", handler.GetMessagesByEmail(api)).Methods(http.MethodGet)

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
