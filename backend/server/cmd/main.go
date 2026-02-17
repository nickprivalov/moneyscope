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

	log.Info("Starting server...")

	// TODO env config this
	dbConfig := db.DbConfig{
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: "password",
		Database: "moneyscope",
		SSLMode:  "disable",
	}

	database, err := db.InitDb(dbConfig)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// TODO move to router
	http.HandleFunc("/dbhealth", func(w http.ResponseWriter, r *http.Request) {
		if err := db.CheckHealth(database); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Database Unoperational"))
			log.Fatalf("Database unoperational: %v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ledger Service is alive!"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.WithField("port", port).Info("Ledger service running.")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
