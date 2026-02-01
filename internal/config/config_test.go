package config

import (
	"testing"
	"time"
)

func TestLoad_Defaults(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Server.Host != "localhost" {
		t.Errorf("Server.Host = %v, want localhost", cfg.Server.Host)
	}
	if cfg.Server.Port != ":8300" {
		t.Errorf("Server.Port = %v, want :8300", cfg.Server.Port)
	}
	if cfg.Server.ReadTimeout != 10*time.Second {
		t.Errorf("Server.ReadTimeout = %v, want 10s", cfg.Server.ReadTimeout)
	}
	if cfg.Database.Port != 5432 {
		t.Errorf("Database.Port = %v, want 5432", cfg.Database.Port)
	}
	if cfg.Database.SSLMode != "require" {
		t.Errorf("Database.SSLMode = %v, want require", cfg.Database.SSLMode)
	}
}

func TestLoad_FromEnv(t *testing.T) {
	t.Setenv("APP_SERVER_HOST", "0.0.0.0")
	t.Setenv("APP_SERVER_PORT", ":9000")
	t.Setenv("APP_DATABASE_HOST", "db.example.com")
	t.Setenv("APP_DATABASE_PORT", "5433")
	t.Setenv("APP_DATABASE_USER", "testuser")
	t.Setenv("APP_DATABASE_PASSWORD", "testpass")
	t.Setenv("APP_DATABASE_DBNAME", "testdb")
	t.Setenv("APP_DATABASE_MAXCONNS", "50")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Server.Host = %v, want 0.0.0.0", cfg.Server.Host)
	}
	if cfg.Server.Port != ":9000" {
		t.Errorf("Server.Port = %v, want :9000", cfg.Server.Port)
	}
	if cfg.Database.Host != "db.example.com" {
		t.Errorf("Database.Host = %v, want db.example.com", cfg.Database.Host)
	}
	if cfg.Database.Port != 5433 {
		t.Errorf("Database.Port = %v, want 5433", cfg.Database.Port)
	}
	if cfg.Database.User != "testuser" {
		t.Errorf("Database.User = %v, want testuser", cfg.Database.User)
	}
	if cfg.Database.MaxConns != 50 {
		t.Errorf("Database.MaxConns = %v, want 50", cfg.Database.MaxConns)
	}
}
