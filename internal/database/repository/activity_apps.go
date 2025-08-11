package repository

import (
	"time"
	
	"shien/internal/utils"
)

// GetAppUsageSummary returns app usage statistics for a given time range
func (r *ActivityRepo) GetAppUsageSummary(from, to time.Time) (map[string]int, error) {
	rows, err := r.conn.Query(`
		SELECT app_name, COUNT(*) as minutes
		FROM activity_logs 
		WHERE recorded_at >= ? 
		  AND recorded_at <= ?
		  AND app_name IS NOT NULL
		GROUP BY app_name
		ORDER BY minutes DESC
	`, utils.ToUTC(from), utils.ToUTC(to))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	usage := make(map[string]int)
	for rows.Next() {
		var appName string
		var minutes int
		err := rows.Scan(&appName, &minutes)
		if err != nil {
			return nil, err
		}
		// Each record represents 5 minutes
		usage[appName] = minutes * 5
	}
	
	return usage, rows.Err()
}

// GetRecentAppActivity returns the most recent app activities
func (r *ActivityRepo) GetRecentAppActivity(limit int) ([]ActivityLog, error) {
	rows, err := r.conn.Query(`
		SELECT id, recorded_at, app_name
		FROM activity_logs 
		WHERE app_name IS NOT NULL
		ORDER BY recorded_at DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var logs []ActivityLog
	for rows.Next() {
		var log ActivityLog
		err := rows.Scan(&log.ID, &log.RecordedAt, &log.AppName)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	
	return logs, rows.Err()
}