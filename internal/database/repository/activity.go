package repository

import (
	"database/sql"
	"time"
	
	"shien/internal/utils"
)

// ActivityLog represents an activity record
type ActivityLog struct {
	ID         int64          `json:"id"`
	RecordedAt utils.UTCTime `json:"recorded_at"`
}

// ActivityRepo implements ActivityRepository
type ActivityRepo struct {
	conn *sql.DB
}

// NewActivityRepo creates a new activity repository
func NewActivityRepo(conn *sql.DB) *ActivityRepo {
	return &ActivityRepo{conn: conn}
}

// RecordActivity records that the app is running at the current time
func (r *ActivityRepo) RecordActivity() error {
	// Round to minute precision
	now := utils.Now().TruncateToMinute()
	
	// Try to insert, ignore if already exists for this minute
	_, err := r.conn.Exec(`
		INSERT OR IGNORE INTO activity_logs (recorded_at) 
		VALUES (?)
	`, now)
	
	return err
}

// GetActivityLogs returns activity logs within a time range
func (r *ActivityRepo) GetActivityLogs(from, to time.Time) ([]ActivityLog, error) {
	rows, err := r.conn.Query(`
		SELECT id, recorded_at 
		FROM activity_logs 
		WHERE recorded_at >= ? 
		  AND recorded_at <= ?
		ORDER BY recorded_at DESC
	`, utils.ToUTC(from), utils.ToUTC(to))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var logs []ActivityLog
	for rows.Next() {
		var log ActivityLog
		err := rows.Scan(&log.ID, &log.RecordedAt)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	
	return logs, rows.Err()
}

// GetActivitySummary returns a summary of activity for a given date
func (r *ActivityRepo) GetActivitySummary(date time.Time) (map[string]interface{}, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	
	// Count total activity records (each represents 5 minutes)
	var count int
	err := r.conn.QueryRow(`
		SELECT COUNT(*) 
		FROM activity_logs 
		WHERE recorded_at >= ? 
		  AND recorded_at < ?
	`, utils.ToUTC(startOfDay), utils.ToUTC(endOfDay)).Scan(&count)
	
	if err != nil {
		return nil, err
	}
	
	return map[string]interface{}{
		"date":           date.Format("2006-01-02"),
		"activity_count": count,
		"minutes_active": count * 5, // Each record represents 5 minutes
	}, nil
}