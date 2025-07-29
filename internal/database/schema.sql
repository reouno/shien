-- Schema definition for Shien database
-- This file is for reference only. Actual schema is managed by migrations.

-- Activity logs table
CREATE TABLE IF NOT EXISTS activity_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    recorded_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_activity_logs_recorded_at 
ON activity_logs(recorded_at);

CREATE UNIQUE INDEX IF NOT EXISTS idx_activity_logs_minute 
ON activity_logs(strftime('%Y-%m-%d %H:%M', recorded_at));

-- Migrations tracking table
CREATE TABLE IF NOT EXISTS migrations (
    version INTEGER PRIMARY KEY,
    description TEXT,
    applied_at DATETIME NOT NULL
);

-- Future tables can be documented here
-- CREATE TABLE users ( ... );
-- CREATE TABLE tasks ( ... );