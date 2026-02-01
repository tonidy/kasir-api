package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

func EnableRLS(db *sql.DB, rlsPath string) error {
	rlsFile := filepath.Join(rlsPath, "enable_rls.sql")

	content, err := os.ReadFile(rlsFile)
	if err != nil {
		return fmt.Errorf("failed to read RLS file: %w", err)
	}

	if _, err := db.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to enable RLS: %w", err)
	}

	return nil
}

func DisableRLS(db *sql.DB, rlsPath string) error {
	rlsFile := filepath.Join(rlsPath, "disable_rls.sql")

	content, err := os.ReadFile(rlsFile)
	if err != nil {
		return fmt.Errorf("failed to read RLS file: %w", err)
	}

	if _, err := db.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to disable RLS: %w", err)
	}

	return nil
}
