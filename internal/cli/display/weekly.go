package display

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"shien/internal/database/repository"
)

// WeeklyReporter handles the display of weekly activity logs
type WeeklyReporter struct{}

// NewWeeklyReporter creates a new weekly reporter
func NewWeeklyReporter() *WeeklyReporter {
	return &WeeklyReporter{}
}

// ShowDailySummary displays activity summary for each day of the last 7 days
func (r *WeeklyReporter) ShowDailySummary(logs []repository.ActivityLog) {
	fmt.Println("Weekly Activity - Daily Summary")
	fmt.Println("================================")
	
	if len(logs) == 0 {
		fmt.Println("No activity logs found for the last 7 days")
		return
	}

	// Group by day
	dailyActivity := make(map[string]int)
	for _, log := range logs {
		day := log.RecordedAt.Format("2006-01-02")
		dailyActivity[day]++
	}

	// Get all 7 days including those with no activity
	now := time.Now()
	days := make([]string, 0, 7)
	for i := 6; i >= 0; i-- {
		day := now.AddDate(0, 0, -i).Format("2006-01-02")
		days = append(days, day)
		if _, exists := dailyActivity[day]; !exists {
			dailyActivity[day] = 0
		}
	}

	// Display daily summary
	fmt.Println("\nLast 7 days:")
	for _, day := range days {
		count := dailyActivity[day]
		minutes := count * 5
		hours := float64(minutes) / 60.0
		// Scale bar: 0-24 hours (1440 minutes) = 0-100% of bar width
		bar := r.makeAbsoluteBar(hours, 24.0, 30)
		
		// Parse day to get weekday name
		t, _ := time.Parse("2006-01-02", day)
		weekday := t.Weekday().String()[:3]
		
		fmt.Printf("%s %s: %s %6.1fh (%3d records)\n", 
			day, weekday, bar, hours, count)
	}
	
	// Total summary
	totalRecords := len(logs)
	totalMinutes := totalRecords * 5
	totalHours := float64(totalMinutes) / 60.0
	fmt.Printf("\nTotal: %.1f hours (%d records)\n", totalHours, totalRecords)
}

// ShowHourlyAverage displays average activity per hour across the last 7 days
func (r *WeeklyReporter) ShowHourlyAverage(logs []repository.ActivityLog) {
	fmt.Println("Weekly Activity - Hourly Average")
	fmt.Println("=================================")
	
	if len(logs) == 0 {
		fmt.Println("No activity logs found for the last 7 days")
		return
	}

	// Group by hour of day
	hourlyActivity := make(map[int]int)
	daysWithActivity := make(map[string]bool)
	
	for _, log := range logs {
		hour := log.RecordedAt.Time.Hour()
		hourlyActivity[hour]++
		day := log.RecordedAt.Format("2006-01-02")
		daysWithActivity[day] = true
	}

	// Calculate number of days (for averaging)
	numDays := 7 // We're looking at last 7 days

	// Display hourly average
	fmt.Println("\nAverage activity by hour (last 7 days):")
	for hour := 0; hour < 24; hour++ {
		count := hourlyActivity[hour]
		avgRecords := float64(count) / float64(numDays)
		avgMinutes := avgRecords * 5
		// Scale bar: 0-60 minutes = 0-100% of bar width
		bar := r.makeAbsoluteBar(avgMinutes, 60.0, 30)
		
		fmt.Printf("%02d:00: %s %5.1fm (avg %.1f records/day)\n", 
			hour, bar, avgMinutes, avgRecords)
	}
	
	// Peak hours
	fmt.Println("\nPeak activity hours:")
	type hourCount struct {
		hour int
		avg  float64
	}
	
	var hours []hourCount
	for h, c := range hourlyActivity {
		hours = append(hours, hourCount{h, float64(c) / float64(numDays)})
	}
	
	sort.Slice(hours, func(i, j int) bool {
		return hours[i].avg > hours[j].avg
	})
	
	// Show top 3 hours
	for i := 0; i < 3 && i < len(hours); i++ {
		if hours[i].avg > 0 {
			fmt.Printf("  %02d:00 - %.1f records/day (%.1f minutes/day)\n", 
				hours[i].hour, hours[i].avg, hours[i].avg*5)
		}
	}
}

// makeAbsoluteBar creates a visual bar with absolute scale
// value: current value
// maxValue: maximum possible value (24 for hours, 60 for minutes)
// maxWidth: maximum width of the bar in characters
func (r *WeeklyReporter) makeAbsoluteBar(value, maxValue float64, maxWidth int) string {
	if value <= 0 {
		return strings.Repeat("·", maxWidth)
	}
	
	// Calculate filled portion based on absolute scale
	ratio := value / maxValue
	if ratio > 1.0 {
		ratio = 1.0 // Cap at 100% if value exceeds max
	}
	
	filled := int(ratio * float64(maxWidth))
	if filled == 0 && value > 0 {
		filled = 1 // Ensure at least one block for non-zero values
	}
	
	bar := strings.Repeat("█", filled)
	empty := strings.Repeat("·", maxWidth-filled)
	return bar + empty
}