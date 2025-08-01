package display

import (
	"fmt"
	"strings"
	"time"

	"shien/internal/database/repository"
)

// ActivityReporter handles the display of activity logs
type ActivityReporter struct{}

// NewActivityReporter creates a new activity reporter
func NewActivityReporter() *ActivityReporter {
	return &ActivityReporter{}
}

// ShowSummary displays the activity summary including hourly breakdown
func (r *ActivityReporter) ShowSummary(logs []repository.ActivityLog) {
	if len(logs) == 0 {
		fmt.Println("No activity logs found for the specified period")
		return
	}

	fmt.Println("Activity Logs")
	fmt.Println("=============")
	fmt.Printf("Total records: %d (≈ %d minutes)\n\n", len(logs), len(logs)*5)

	r.showHourlyBreakdown(logs)
}

// showHourlyBreakdown displays activity grouped by hour with visual bars
func (r *ActivityReporter) showHourlyBreakdown(logs []repository.ActivityLog) {
	// Group by hour for display
	hourlyCount := make(map[string]int)
	for _, log := range logs {
		hour := log.RecordedAt.Format("2006-01-02 15:00")
		hourlyCount[hour]++
	}

	// Find the time range to fill in missing hours
	if len(logs) == 0 {
		fmt.Println("No activity recorded.")
		return
	}

	// Find min and max times
	minTime := logs[0].RecordedAt.Time
	maxTime := logs[0].RecordedAt.Time
	for _, log := range logs {
		if log.RecordedAt.Time.Before(minTime) {
			minTime = log.RecordedAt.Time
		}
		if log.RecordedAt.Time.After(maxTime) {
			maxTime = log.RecordedAt.Time
		}
	}

	// Round down to start of hour for minTime
	startHour := time.Date(minTime.Year(), minTime.Month(), minTime.Day(), minTime.Hour(), 0, 0, 0, minTime.Location())
	endHour := time.Date(maxTime.Year(), maxTime.Month(), maxTime.Day(), maxTime.Hour(), 0, 0, 0, maxTime.Location())

	// Generate all hours in the range
	var hours []string
	for h := startHour; !h.After(endHour); h = h.Add(time.Hour) {
		hourStr := h.Format("2006-01-02 15:00")
		hours = append(hours, hourStr)
		// Initialize zero count for hours without activity
		if _, exists := hourlyCount[hourStr]; !exists {
			hourlyCount[hourStr] = 0
		}
	}

	fmt.Println("Activity by hour:")
	for _, hour := range hours {
		count := hourlyCount[hour]
		bar := r.makeBar(count)
		fmt.Printf("%s: %s (%d)\n", hour, bar, count*5)
	}
}

// makeBar creates a visual bar representation
func (r *ActivityReporter) makeBar(count int) string {
	if count == 0 {
		return "-"
	}
	return strings.Repeat("█", count)
}