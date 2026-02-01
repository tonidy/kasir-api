package database

import (
	"testing"

	"kasir-api/internal/config"
)

func TestBuildDSN(t *testing.T) {
	cfg := config.DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "testuser",
		Password: "testpass",
		DBName:   "testdb",
		SSLMode:  "disable",
		MaxConns: 10,
		MinConns: 2,
	}

	db, err := NewPool(cfg)
	if err != nil {
		t.Fatalf("NewPool() error = %v", err)
	}
	defer db.Close()

	if db.DB == nil {
		t.Error("DB is nil")
	}
}
