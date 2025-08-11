package migrations

import (
	"database/sql"
)

// Migration003_AddAppNameToActivity adds app_name column to activity_logs
var Migration003_AddAppNameToActivity = Migration{
	Version:     3,
	Description: "Add app_name column to activity_logs",
	Up: func(tx *sql.Tx) error {
		// Add app_name column to track which application was being used
		if _, err := tx.Exec(`
			ALTER TABLE activity_logs 
			ADD COLUMN app_name TEXT
		`); err != nil {
			return err
		}
		
		// Add index for efficient queries by app_name
		if _, err := tx.Exec(`
			CREATE INDEX IF NOT EXISTS idx_activity_logs_app_name 
			ON activity_logs(app_name)
		`); err != nil {
			return err
		}
		
		// Add composite index for date and app queries
		if _, err := tx.Exec(`
			CREATE INDEX IF NOT EXISTS idx_activity_logs_date_app 
			ON activity_logs(recorded_at, app_name)
		`); err != nil {
			return err
		}
		
		return nil
	},
}