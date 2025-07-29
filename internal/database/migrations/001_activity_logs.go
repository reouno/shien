package migrations

import (
	"database/sql"
)

// Migration001ActivityLogs creates the activity logs table
var Migration001ActivityLogs = Migration{
	Version:     1,
	Description: "Create activity logs table",
	Up: func(tx *sql.Tx) error {
		// Activity logs table for tracking app is running
		if _, err := tx.Exec(`
			CREATE TABLE IF NOT EXISTS activity_logs (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				recorded_at DATETIME NOT NULL
			)
		`); err != nil {
			return err
		}
		
		// Index for querying by date
		if _, err := tx.Exec(`
			CREATE INDEX IF NOT EXISTS idx_activity_logs_recorded_at 
			ON activity_logs(recorded_at)
		`); err != nil {
			return err
		}
		
		// Create unique index to prevent duplicate entries for the same minute
		if _, err := tx.Exec(`
			CREATE UNIQUE INDEX IF NOT EXISTS idx_activity_logs_minute 
			ON activity_logs(strftime('%Y-%m-%d %H:%M', recorded_at))
		`); err != nil {
			return err
		}
		
		return nil
	},
}