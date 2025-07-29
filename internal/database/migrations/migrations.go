package migrations

import (
	"database/sql"
)

// Migration represents a database migration
type Migration struct {
	Version     int
	Description string
	Up          func(*sql.Tx) error
}

// All returns all migrations in order
func All() []Migration {
	return []Migration{
		Migration001ActivityLogs,
		// Future migrations will be added here:
		// Migration002UserTable,
		// Migration003AddFieldToActivityLogs,
	}
}