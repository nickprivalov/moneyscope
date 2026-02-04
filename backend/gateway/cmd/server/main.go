package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	ingestv1 "github.com/nickprivalov/moneyscope/backend/gen/go/ingest/v1"
	ledgerv1 "github.com/nickprivalov/moneyscope/backend/gen/go/ledger/v1"
)

func main() {
	log.SetOutput(os.Stdout)
	log.Println("Starting Gateway Service...")

	// 1. Connect to Microservices
	//ledgerClient := connectLedger()
	ingestClient := connectIngest()

	// 2. Setup HTTP Router (Gin)
	r := gin.Default()

	// CORS Middleware (Crucial for React Frontend)
	r.Use(func(c *gin.Context) {

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 3. Define Routes
	api := r.Group("/api")
	{
		api.POST("/upload", func(c *gin.Context) {
			handleUpload(c, ingestClient)
		})

		// TODO: Add Transaction Listing
		// api.GET("/transactions", ...)
	}

	// 4. Start Server
	log.Println("Gateway listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run gateway: %v", err)
	}
}

// --- Handlers ---

func handleUpload(c *gin.Context, ingestClient ingestv1.IngestServiceClient) {
	// 1. Get File from Multipart Form
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer file.Close()

	// 2. Start gRPC Stream
	stream, err := ingestClient.UploadStatement(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to ingestion service"})
		return
	}

	// 3. Send Metadata (First Message)
	err = stream.Send(&ingestv1.UploadStatementRequest{
		Payload: &ingestv1.UploadStatementRequest_Metadata{
			Metadata: &ingestv1.UploadMetadata{
				UserId:    "test-user-id",    // Hardcoded for now (Auth comes later)
				AccountId: "test-account-id", // Hardcoded for now
				FileName:  fileHeader.Filename,
			},
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send metadata"})
		return
	}

	// 4. Stream File Chunks
	buffer := make([]byte, 1024*64) // 64KB chunks
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
			return
		}

		err = stream.Send(&ingestv1.UploadStatementRequest{
			Payload: &ingestv1.UploadStatementRequest_Chunk{
				Chunk: buffer[:n],
			},
		})
		if err != nil {
			log.Printf("Error sending chunk: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Stream interrupted"})
			return
		}
	}

	// 5. Close Stream & Get Response
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Printf("Error closing stream: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Processing failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// --- Connection Helpers ---

func connectLedger() ledgerv1.LedgerServiceClient {
	host := os.Getenv("LEDGER_HOST")
	if host == "" {
		host = "localhost"
	}
	conn, err := grpc.NewClient(host+":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Ledger: %v", err)
	}
	return ledgerv1.NewLedgerServiceClient(conn)
}

func connectIngest() ingestv1.IngestServiceClient {
	host := os.Getenv("INGEST_HOST")
	if host == "" {
		host = "localhost"
	}
	conn, err := grpc.NewClient(host+":50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Ingest: %v", err)
	}
	return ingestv1.NewIngestServiceClient(conn)
}
