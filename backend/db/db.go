package db

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type DbConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	SSLMode  string
}

func InitDb(cfg DbConfig) (*sqlx.DB, error) {
	log.WithFields(log.Fields{
		"host":     cfg.Host,
		"port":     cfg.Port,
		"user":     cfg.Username,
		"database": cfg.Database,
		"ssl":      cfg.SSLMode,
	})

	// connection string
	connectString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database, cfg.SSLMode,
	)

	// connect to DB
	db, err := sqlx.Connect("postgres", connectString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Successfully connected to database")

	return db, nil
}

func RunMigrations(db *sqlx.DB) error {
	panic("not implemented: apply SQL migration files")
}

func WithTransaction(db *sqlx.DB, f func(tx *sql.Tx) error) error {
	panic("not implemented: begin/commit/rollback DB transaction")
}

func CheckHealth(db *sqlx.DB) error {
	panic("not implemented: pinging DB and getting health status")
}
