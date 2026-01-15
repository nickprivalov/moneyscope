package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"

	"github.com/nickprivalov/moneyscope/backend/common/db"
	ledgerv1 "github.com/nickprivalov/moneyscope/backend/gen/go/ledger/v1"
	"github.com/nickprivalov/moneyscope/backend/ledger/internal/models"
	"github.com/nickprivalov/moneyscope/backend/ledger/internal/service"
)

func main() {
	log.SetOutput(os.Stdout)
	log.Println("Starting Ledger Service...")

	// 1. Config (Ideally load from env vars, hardcoded for now for MVP)
	dbConfig := db.DbConfig{
		Host:     "localhost", // Use "postgres" if running inside docker network
		Port:     "5432",
		User:     "moneyscope_user",
		Password: "moneyscope_password",
		DBName:   "moneyscope_db",
		SSLMode:  "disable",
	}

	// Override host if running in Docker (set ENV var DB_HOST=postgres)
	if host := os.Getenv("DB_HOST"); host != "" {
		dbConfig.Host = host
	}

	// 2. Database Connection
	gormDB, err := db.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 3. Migrations
	log.Println("Running Database Migrations...")
	if err := runMigrations(gormDB); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// 4. gRPC Server Setup
	lis, err := net.Listen("tcp", ":50051") // Ledger Service Port
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	ledgerService := service.NewLedgerService(gormDB)

	ledgerv1.RegisterLedgerServiceServer(grpcServer, ledgerService)

	// Enable Reflection (allows tools like Postman/BloomRPC to inspect the schema)
	reflection.Register(grpcServer)

	log.Printf("Ledger Service listening on %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func runMigrations(db *gorm.DB) error {
	// Ensure the schema exists
	if err := db.Exec("CREATE SCHEMA IF NOT EXISTS ledger").Error; err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	// Auto-migrate the tables
	if err := db.AutoMigrate(&models.Transaction{}); err != nil {
		return fmt.Errorf("failed to auto-migrate: %w", err)
	}
	return nil
}
