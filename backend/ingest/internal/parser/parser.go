package parser

import (
	"context"
	"encoding/csv"
	"io"
	"log"
	"strconv"
	"time"

	commonv1 "github.com/nickprivalov/moneyscope/backend/gen/go/common/v1"
	ledgerv1 "github.com/nickprivalov/moneyscope/backend/gen/go/ledger/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ProcessCSV reads from the reader, parses lines, and sends batches to Ledger Service.
func ProcessCSV(ctx context.Context, r io.Reader, client ledgerv1.LedgerServiceClient, userID, accountID string) (int, error) {
	csvReader := csv.NewReader(r)

	// Skip Header (Assumption: Row 1 is header)
	if _, err := csvReader.Read(); err != nil {
		return 0, err // Empty file or read error
	}

	var batch []*ledgerv1.CreateTransactionRequest
	totalProcessed := 0
	batchSize := 50

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading CSV line: %v", err)
			continue // Skip bad lines
		}

		// PARSING LOGIC (MVP Assumption: Date, Description, Amount)
		// e.g., 2025-01-01, Starbucks, -5.50
		if len(record) < 3 {
			continue
		}

		// 1. Parse Date
		date, err := time.Parse("2006-01-02", record[0])
		if err != nil {
			log.Printf("Skipping invalid date: %s", record[0])
			continue
		}

		// 2. Parse Amount (Naive implementation for float string)
		// In production, we'd use a robust money parser library.
		amountFloat, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			log.Printf("Skipping invalid amount: %s", record[2])
			continue
		}

		units := int64(amountFloat)
		nanos := int32((amountFloat - float64(units)) * 1_000_000_000)

		// Determine Type
		txType := ledgerv1.TransactionType_EXPENSE
		if amountFloat > 0 {
			txType = ledgerv1.TransactionType_INCOME
		}

		// 3. Build Request
		req := &ledgerv1.CreateTransactionRequest{
			UserId:      userID,
			AccountId:   accountID,
			Date:        timestamppb.New(date),
			Description: record[1],
			Type:        txType,
			Amount: &commonv1.Money{
				CurrencyCode: "USD",
				Units:        units,
				Nanos:        nanos,
			},
		}

		batch = append(batch, req)

		// 4. Send Batch if full
		if len(batch) >= batchSize {
			if err := sendBatch(ctx, client, batch); err != nil {
				log.Printf("Failed to send batch: %v", err)
			} else {
				totalProcessed += len(batch)
			}
			batch = nil // Reset
		}
	}

	// Send remaining
	if len(batch) > 0 {
		if err := sendBatch(ctx, client, batch); err == nil {
			totalProcessed += len(batch)
		}
	}

	return totalProcessed, nil
}

func sendBatch(ctx context.Context, client ledgerv1.LedgerServiceClient, txs []*ledgerv1.CreateTransactionRequest) error {
	// We need to add CreateTransactionsBatch to LedgerService first!
	// For MVP, loop and call single Create. Optimally, use Batch RPC.

	// Since we defined CreateTransactionsBatch in proto, let's use it.
	_, err := client.CreateTransactionsBatch(ctx, &ledgerv1.CreateTransactionsBatchRequest{
		Transactions: txs,
	})
	return err
}
