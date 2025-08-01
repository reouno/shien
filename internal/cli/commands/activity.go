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

// ActivityCommand handles activity log display
type ActivityCommand struct{}

// NewActivityCommand creates a new activity command
func NewActivityCommand() *ActivityCommand {
	return &ActivityCommand{}
}

// Name returns the command name
func (c *ActivityCommand) Name() string {
	return "activity"
}

// Description returns the command description
func (c *ActivityCommand) Description() string {
	return "Show activity logs"
}

// Usage returns the command usage
func (c *ActivityCommand) Usage() string {
	return `activity [options]
    -from <date>      Start date (YYYY-MM-DD)
    -to <date>        End date (YYYY-MM-DD)
    -today            Show today's activity`
}

// Execute runs the activity command
func (c *ActivityCommand) Execute(client *rpc.Client, args []string) error {
	flags := flag.NewFlagSet("activity", flag.ExitOnError)
	from := flags.String("from", "", "Start date (YYYY-MM-DD)")
	to := flags.String("to", "", "End date (YYYY-MM-DD)")
	today := flags.Bool("today", false, "Show today's activity")

	if err := flags.Parse(args); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	params := make(map[string]interface{})

	if *today {
		now := time.Now()
		params["from"] = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Format(time.RFC3339)
		params["to"] = now.Format(time.RFC3339)
	} else {
		if *from != "" {
			t, err := time.Parse("2006-01-02", *from)
			if err != nil {
				return fmt.Errorf("invalid from date: %w", err)
			}
			params["from"] = t.Format(time.RFC3339)
		}

		if *to != "" {
			t, err := time.Parse("2006-01-02", *to)
			if err != nil {
				return fmt.Errorf("invalid to date: %w", err)
			}
			// Set to end of day
			t = t.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			params["to"] = t.Format(time.RFC3339)
		}
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

	// Display the activity report
	reporter := display.NewActivityReporter()
	reporter.ShowSummary(logs)

	return nil
}