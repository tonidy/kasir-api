package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"kasir-api/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	*sql.DB
}

func NewPool(cfg config.DatabaseConfig) (*DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s default_query_exec_mode=simple_protocol",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	// Open database connection
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxConns)
	db.SetMaxIdleConns(cfg.MinConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Minute)
	db.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleTime) * time.Minute)

	return &DB{db}, nil
}

func (db *DB) Ping(ctx context.Context) error {
	return db.PingContext(ctx)
}
