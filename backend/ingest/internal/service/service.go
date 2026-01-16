package service

import (
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ingestv1 "github.com/nickprivalov/moneyscope/backend/gen/go/ingest/v1"
	ledgerv1 "github.com/nickprivalov/moneyscope/backend/gen/go/ledger/v1"
	"github.com/nickprivalov/moneyscope/backend/ingest/internal/parser"
)

type IngestService struct {
	ingestv1.UnimplementedIngestServiceServer
	ledgerClient ledgerv1.LedgerServiceClient
}

func NewIngestService(ledgerClient ledgerv1.LedgerServiceClient) *IngestService {
	return &IngestService{
		ledgerClient: ledgerClient,
	}
}

func (s *IngestService) UploadStatement(stream ingestv1.IngestService_UploadStatementServer) error {
	// 1. Receive the first message (Metadata)
	req, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.Unknown, "cannot receive metadata: %v", err)
	}

	metadata := req.GetMetadata()
	if metadata == nil {
		return status.Error(codes.InvalidArgument, "first message must be metadata")
	}

	log.Printf("Starting upload for User: %s, File: %s", metadata.UserId, metadata.FileName)

	// 2. Create a pipe to connect the gRPC stream to the CSV reader
	reader, writer := io.Pipe()

	// 3. Launch a goroutine to read from gRPC and write to the pipe
	go func() {
		defer func(writer *io.PipeWriter) {
			err := writer.Close()
			if err != nil {
				log.Printf("Error closing writer: %v", err)
			}
		}(writer)
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Printf("Error receiving chunk: %v", err)
				_ = writer.CloseWithError(err)
				return
			}

			chunk := req.GetChunk()
			if chunk == nil {
				continue // skip empty
			}

			if _, err := writer.Write(chunk); err != nil {
				log.Printf("Error writing to pipe: %v", err)
				return
			}
		}
	}()

	// 4. Process the stream using the Parser
	// This blocks until the pipe is closed (EOF from stream)
	processedCount, err := parser.ProcessCSV(stream.Context(), reader, s.ledgerClient, metadata.UserId, metadata.AccountId)
	if err != nil {
		return status.Errorf(codes.Internal, "processing failed: %v", err)
	}

	// 5. Send Response
	return stream.SendAndClose(&ingestv1.UploadStatementResponse{
		TransactionsProcessed: int32(processedCount),
		Message:               fmt.Sprintf("Successfully processed %d transactions", processedCount),
	})
}
