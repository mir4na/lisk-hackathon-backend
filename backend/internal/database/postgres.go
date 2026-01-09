package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/vessel/backend/internal/config"
)

func NewPostgresConnection(cfg *config.Config) (*sql.DB, error) {
	connStr := cfg.DatabaseURL
	if connStr == "" {
		connStr = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.PostgresHost,
			cfg.PostgresPort,
			cfg.PostgresUser,
			cfg.PostgresPass,
			cfg.PostgresDB,
		)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Connected to PostgreSQL database")
	return db, nil
}

type TxFn func(*sql.Tx) error

func WithTransaction(db *sql.DB, fn TxFn) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
