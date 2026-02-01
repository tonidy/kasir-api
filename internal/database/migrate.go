package database

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func RunMigrations(db *sql.DB, migrationsPath string) error {
	goose.SetDialect("postgres")

	if err := goose.Up(db, migrationsPath); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func ResetMigrations(db *sql.DB, migrationsPath string) error {
	goose.SetDialect("postgres")

	if err := goose.Reset(db, migrationsPath); err != nil {
		return fmt.Errorf("failed to reset migrations: %w", err)
	}

	return nil
}
