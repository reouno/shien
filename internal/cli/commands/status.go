package commands

import (
	"fmt"
	"time"

	"shien/internal/rpc"
)

// StatusCommand handles status display
type StatusCommand struct{}

// NewStatusCommand creates a new status command
func NewStatusCommand() *StatusCommand {
	return &StatusCommand{}
}

// Name returns the command name
func (c *StatusCommand) Name() string {
	return "status"
}

// Description returns the command description
func (c *StatusCommand) Description() string {
	return "Show daemon status"
}

// Usage returns the command usage
func (c *StatusCommand) Usage() string {
	return "status"
}

// Execute runs the status command
func (c *StatusCommand) Execute(client *rpc.Client, args []string) error {
	status, err := client.GetStatus()
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	fmt.Println("Shien Service Status")
	fmt.Println("==================")
	fmt.Printf("Running: %v\n", status.Running)
	fmt.Printf("Started: %s\n", status.StartedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Uptime:  %s\n", time.Since(status.StartedAt).Round(time.Second))
	fmt.Printf("Version: %s\n", status.Version)

	return nil
}