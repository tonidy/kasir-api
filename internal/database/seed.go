package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

func RunSeeds(db *sql.DB, seedsPath string) error {
	seedFile := filepath.Join(seedsPath, "seed.sql")

	content, err := os.ReadFile(seedFile)
	if err != nil {
		return fmt.Errorf("failed to read seed file: %w", err)
	}

	if _, err := db.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute seed: %w", err)
	}

	return nil
}
