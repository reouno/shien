package migrations

import (
	"database/sql"
)

// Migration002_Gamification adds gamification tables
var Migration002_Gamification = Migration{
	Version:     2,
	Description: "Add gamification tables",
	Up: func(tx *sql.Tx) error {
		// User gamification status table
		if _, err := tx.Exec(`
			CREATE TABLE IF NOT EXISTS user_status (
			user_id TEXT PRIMARY KEY,
			level INTEGER NOT NULL DEFAULT 1,
			experience INTEGER NOT NULL DEFAULT 0,
			total_exp INTEGER NOT NULL DEFAULT 0,
			
			-- Work-related attributes (no upper limit)
			focus INTEGER NOT NULL DEFAULT 50,
			productivity INTEGER NOT NULL DEFAULT 50,
			creativity INTEGER NOT NULL DEFAULT 50,
			stamina INTEGER NOT NULL DEFAULT 100,
			knowledge INTEGER NOT NULL DEFAULT 10,
			collaboration INTEGER NOT NULL DEFAULT 30,
			
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
		`); err != nil {
			return err
		}

		// Attribute modifiers table (for temporary/permanent effects)
		if _, err := tx.Exec(`
			CREATE TABLE IF NOT EXISTS attribute_modifiers (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			attribute TEXT NOT NULL, -- 'focus', 'productivity', etc.
			value INTEGER NOT NULL, -- can be positive or negative
			reason TEXT NOT NULL,
			expires_at DATETIME, -- NULL for permanent modifiers
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			
			FOREIGN KEY (user_id) REFERENCES user_status(user_id)
		)
		`); err != nil {
			return err
		}

		// Index for efficient queries
		if _, err := tx.Exec(`
			CREATE INDEX IF NOT EXISTS idx_modifiers_user_expires 
			ON attribute_modifiers(user_id, expires_at)
		`); err != nil {
			return err
		}

		// Achievement tracking table (for future use)
		if _, err := tx.Exec(`
			CREATE TABLE IF NOT EXISTS achievements (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			achievement_type TEXT NOT NULL,
			unlocked_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			
			FOREIGN KEY (user_id) REFERENCES user_status(user_id),
			UNIQUE(user_id, achievement_type)
		)
		`); err != nil {
			return err
		}

		return nil
	},
}