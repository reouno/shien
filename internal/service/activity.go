package service

import (
	"time"
	
	"shien/internal/database/repository"
)

// ActivityService handles business logic for activity tracking
type ActivityService struct {
	repo *repository.ActivityRepo
}

// NewActivityService creates a new activity service
func NewActivityService(repo *repository.ActivityRepo) *ActivityService {
	return &ActivityService{repo: repo}
}

// RecordActivity records current activity
func (s *ActivityService) RecordActivity() error {
	return s.repo.RecordActivity()
}

// GetActivityLogs retrieves activity logs for a date range
func (s *ActivityService) GetActivityLogs(from, to time.Time) ([]repository.ActivityLog, error) {
	// Business logic validation
	if from.After(to) {
		from, to = to, from
	}
	
	// Default to last 24 hours if not specified
	if from.IsZero() {
		from = time.Now().Add(-24 * time.Hour)
	}
	if to.IsZero() {
		to = time.Now()
	}
	
	return s.repo.GetActivityLogs(from, to)
}

// GetActivitySummary gets activity summary for a date
func (s *ActivityService) GetActivitySummary(date time.Time) (map[string]interface{}, error) {
	return s.repo.GetActivitySummary(date)
}

// GetDailyStats returns statistics for today
func (s *ActivityService) GetDailyStats() (map[string]interface{}, error) {
	today := time.Now()
	logs, err := s.repo.GetActivityLogs(
		time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location()),
		today,
	)
	if err != nil {
		return nil, err
	}
	
	return map[string]interface{}{
		"date":          today.Format("2006-01-02"),
		"record_count":  len(logs),
		"minutes_active": len(logs) * 5,
		"hours_active":  float64(len(logs)*5) / 60.0,
	}, nil
}