package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	
	"shien/internal/database/migrations"
)

// Migrate runs all pending database migrations
func (db *DB) Migrate() error {
	// Create migrations table
	if err := db.createMigrationsTable(); err != nil {
		return err
	}
	
	// Get current version
	currentVersion, err := db.getCurrentVersion()
	if err != nil {
		return err
	}
	
	// Run migrations
	for _, migration := range migrations.All() {
		if migration.Version <= currentVersion {
			continue
		}
		
		log.Printf("Running migration %d: %s", migration.Version, migration.Description)
		
		if err := db.Transaction(func(tx *sql.Tx) error {
			// Run migration
			if err := migration.Up(tx); err != nil {
				return fmt.Errorf("migration %d failed: %w", migration.Version, err)
			}
			
			// Update version
			_, err := tx.Exec(
				"INSERT INTO migrations (version, description, applied_at) VALUES (?, ?, ?)",
				migration.Version, migration.Description, time.Now(),
			)
			return err
		}); err != nil {
			return err
		}
		
		log.Printf("Migration %d completed", migration.Version)
	}
	
	return nil
}

// createMigrationsTable creates the migrations tracking table
func (db *DB) createMigrationsTable() error {
	_, err := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			version INTEGER PRIMARY KEY,
			description TEXT,
			applied_at DATETIME NOT NULL
		)
	`)
	return err
}

// getCurrentVersion returns the current migration version
func (db *DB) getCurrentVersion() (int, error) {
	var version sql.NullInt64
	err := db.conn.QueryRow(
		"SELECT MAX(version) FROM migrations",
	).Scan(&version)
	
	if err != nil {
		return 0, err
	}
	
	if !version.Valid {
		return 0, nil
	}
	
	return int(version.Int64), nil
}