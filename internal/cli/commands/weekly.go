package commands

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"shien/internal/cli/display"
	"shien/internal/database/repository"
	"shien/internal/rpc"
)

// WeeklyCommand handles weekly activity display
type WeeklyCommand struct{}

// NewWeeklyCommand creates a new weekly command
func NewWeeklyCommand() *WeeklyCommand {
	return &WeeklyCommand{}
}

// Name returns the command name
func (c *WeeklyCommand) Name() string {
	return "weekly"
}

// Description returns the command description
func (c *WeeklyCommand) Description() string {
	return "Show weekly activity summary"
}

// Usage returns the command usage
func (c *WeeklyCommand) Usage() string {
	return `weekly [options]
    -daily            Show daily summary (default)
    -hourly           Show hourly average`
}

// Execute runs the weekly command
func (c *WeeklyCommand) Execute(client *rpc.Client, args []string) error {
	flags := flag.NewFlagSet("weekly", flag.ExitOnError)
	daily := flags.Bool("daily", false, "Show daily summary")
	hourly := flags.Bool("hourly", false, "Show hourly average")

	if err := flags.Parse(args); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	// Default to daily view if no flag specified
	if !*daily && !*hourly {
		*daily = true
	}

	// Get last 7 days of data
	now := time.Now()
	weekAgo := now.AddDate(0, 0, -7)
	
	params := map[string]interface{}{
		"from": weekAgo.Format(time.RFC3339),
		"to":   now.Format(time.RFC3339),
	}

	resp, err := client.Call(rpc.MethodGetActivityLogs, params)
	if err != nil {
		return fmt.Errorf("failed to get activity logs: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("error: %s", resp.Error)
	}

	// Convert response data to activity logs
	data, err := json.Marshal(resp.Data)
	if err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	var logs []repository.ActivityLog
	if err := json.Unmarshal(data, &logs); err != nil {
		return fmt.Errorf("failed to parse activity logs: %w", err)
	}

	// Display the weekly report
	reporter := display.NewWeeklyReporter()
	if *hourly {
		reporter.ShowHourlyAverage(logs)
	} else {
		reporter.ShowDailySummary(logs)
	}

	return nil
}