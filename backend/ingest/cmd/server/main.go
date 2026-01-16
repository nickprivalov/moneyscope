package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	ingestv1 "github.com/nickprivalov/moneyscope/backend/gen/go/ingest/v1"
	ledgerv1 "github.com/nickprivalov/moneyscope/backend/gen/go/ledger/v1"
	"github.com/nickprivalov/moneyscope/backend/ingest/internal/service"
)

func main() {
	log.SetOutput(os.Stdout)
	log.Println("Starting Ingest Service...")

	ledgerAddr := "localhost:50051"
	if host := os.Getenv("LEDGER_HOST"); host != "" {
		ledgerAddr = host + ":50051"
	}

	// 1. Connect to Ledger Service
	// We use insecure credentials for internal microservice communication
	conn, err := grpc.NewClient(ledgerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to ledger: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Printf("Error closing db connection: %v", err)
		}
	}(conn)
	ledgerClient := ledgerv1.NewLedgerServiceClient(conn)

	// 2. Setup Ingest Server
	lis, err := net.Listen("tcp", ":50052") // Port 50052
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	ingestService := service.NewIngestService(ledgerClient)

	ingestv1.RegisterIngestServiceServer(grpcServer, ingestService)
	reflection.Register(grpcServer)

	log.Printf("Ingest Service listening on %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
