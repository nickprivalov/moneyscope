package main

import (
	"backend/common"
	"backend/db"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	common.InitLogger()

	log.Info("Starting main server...")

	// TODO env config this
	dbConfig := db.DbConfig{
		Host:     "localhost",
		Port:     "5432",
		Username: "moneyscope_user",
		Password: "moneyscope_password_1234",
		Database: "moneyscope_db",
		SSLMode:  "disable",
	}

	// DB connect
	database, err := db.InitDb(dbConfig)
	if err != nil {
		log.Fatalf("Failed to initialize database connection: %v", err)
	}
	defer database.Close()

	// TODO move to router
	// health check endpoint for testing/running purposes
	http.HandleFunc("/dbhealth", func(w http.ResponseWriter, r *http.Request) {
		if err := db.CheckHealth(database); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Database Unoperational"))
			log.Fatalf("Database unoperational: %v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Main server is alive!"))
	})

	// TODO env config this
	// get running port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// graceful shutdown handling
	// NOTE unneeded with httpserve
	//stopSignal := make(chan os.Signal, 1)
	//signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)

	// serve HTTP
	log.WithField("port", port).Info("Main server is now running.")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// shutdown detection
	//<-stopSignal
	//log.Info("Shutting down server...")

}
