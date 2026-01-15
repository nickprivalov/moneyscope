package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	commonv1 "github.com/nickprivalov/moneyscope/backend/gen/go/common/v1"
	ledgerv1 "github.com/nickprivalov/moneyscope/backend/gen/go/ledger/v1"
	"github.com/nickprivalov/moneyscope/backend/ledger/internal/models"
)

// LedgerService implements the gRPC LedgerServiceServer interface.
type LedgerService struct {
	ledgerv1.UnimplementedLedgerServiceServer
	db *gorm.DB
}

// NewLedgerService creates a new instance of the LedgerService.
func NewLedgerService(db *gorm.DB) *LedgerService {
	return &LedgerService{db: db}
}

// CreateTransaction handles the creation of a single transaction.
func (s *LedgerService) CreateTransaction(ctx context.Context, req *ledgerv1.CreateTransactionRequest) (*ledgerv1.CreateTransactionResponse, error) {
	// 1. Validation (Basic)
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.Amount == nil {
		return nil, status.Error(codes.InvalidArgument, "amount is required")
	}

	// 2. Map Proto -> DB Model
	// We generate the UUID here in Go if we don't want to rely on the DB extension being present yet,
	// or we can let the DB handle it. Let's generate it here for safety.
	newID := uuid.New().String()

	txModel := models.Transaction{
		ID:           newID,
		UserID:       req.UserId,
		AccountID:    req.AccountId,
		Date:         req.Date.AsTime(),
		CurrencyCode: req.Amount.CurrencyCode,
		AmountUnits:  req.Amount.Units,
		AmountNanos:  req.Amount.Nanos,
		Description:  req.Description,
		CategoryID:   req.CategoryId,
		Type:         int32(req.Type),
		// Tags:         req.Tags, // TODO: Handle Tags mapping if needed (requires type conversion)
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 3. Persist to DB
	if err := s.db.Create(&txModel).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create transaction: %v", err)
	}

	// 4. Map DB Model -> Proto Response
	return &ledgerv1.CreateTransactionResponse{
		Transaction: &ledgerv1.Transaction{
			Id:        txModel.ID,
			UserId:    txModel.UserID,
			AccountId: txModel.AccountID,
			Date:      timestamppb.New(txModel.Date),
			Amount: &commonv1.Money{
				CurrencyCode: txModel.CurrencyCode,
				Units:        txModel.AmountUnits,
				Nanos:        txModel.AmountNanos,
			},
			Description: txModel.Description,
			CategoryId:  txModel.CategoryID,
			Type:        ledgerv1.TransactionType(txModel.Type),
			Tags:        nil, // Populate if we saved tags
			CreatedAt:   timestamppb.New(txModel.CreatedAt),
			UpdatedAt:   timestamppb.New(txModel.UpdatedAt),
		},
	}, nil
}
