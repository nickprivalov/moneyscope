package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Transaction represents the database schema for a financial transaction.
// It belongs to the "ledger" schema.
type Transaction struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    string    `gorm:"index;not null"` // Partition key usually
	AccountID string    `gorm:"index;not null"`
	Date      time.Time `gorm:"index;not null"`

	// Composite Money representation
	CurrencyCode string `gorm:"size:3;not null"`
	AmountUnits  int64  `gorm:"not null"`
	AmountNanos  int32  `gorm:"not null"`

	Description string
	CategoryID  string `gorm:"index"`

	// Enum stored as integer (matches Proto enum)
	Type int32 `gorm:"not null"`

	// Postgres Array type for tags
	Tags pq.StringArray `gorm:"type:text[]"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName overrides the table name to include the schema.
// We want to keep ledger data in a 'ledger' schema for logical separation.
func (Transaction) TableName() string {
	return "ledger.transactions"
}
