package repository

import (
	"database/sql"
	"time"
)

// ActivityLog represents an activity record
type ActivityLog struct {
	ID         int64     `json:"id"`
	RecordedAt time.Time `json:"recorded_at"`
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
	now := time.Now().Truncate(time.Minute)
	
	// Try to insert, ignore if already exists for this minute
	_, err := r.conn.Exec(`
		INSERT OR IGNORE INTO activity_logs (recorded_at) 
		VALUES (datetime(?, 'localtime'))
	`, now.Format("2006-01-02 15:04:00"))
	
	return err
}

// GetActivityLogs returns activity logs within a time range
func (r *ActivityRepo) GetActivityLogs(from, to time.Time) ([]ActivityLog, error) {
	rows, err := r.conn.Query(`
		SELECT id, recorded_at 
		FROM activity_logs 
		WHERE recorded_at >= datetime(?, 'localtime') 
		  AND recorded_at <= datetime(?, 'localtime')
		ORDER BY recorded_at DESC
	`, from.Format("2006-01-02 15:04:00"), to.Format("2006-01-02 15:04:00"))
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
		WHERE recorded_at >= datetime(?, 'localtime') 
		  AND recorded_at < datetime(?, 'localtime')
	`, startOfDay.Format("2006-01-02 15:04:00"), endOfDay.Format("2006-01-02 15:04:00")).Scan(&count)
	
	if err != nil {
		return nil, err
	}
	
	return map[string]interface{}{
		"date":           date.Format("2006-01-02"),
		"activity_count": count,
		"minutes_active": count * 5, // Each record represents 5 minutes
	}, nil
}