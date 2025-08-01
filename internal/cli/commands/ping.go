package commands

import (
	"fmt"
	"os"

	"shien/internal/rpc"
)

// PingCommand checks if daemon is running
type PingCommand struct{}

// NewPingCommand creates a new ping command
func NewPingCommand() *PingCommand {
	return &PingCommand{}
}

// Name returns the command name
func (c *PingCommand) Name() string {
	return "ping"
}

// Description returns the command description
func (c *PingCommand) Description() string {
	return "Check if daemon is running"
}

// Usage returns the command usage
func (c *PingCommand) Usage() string {
	return "ping"
}

// Execute runs the ping command
func (c *PingCommand) Execute(client *rpc.Client, args []string) error {
	if err := client.Ping(); err != nil {
		fmt.Println("❌ Daemon is not running")
		os.Exit(1)
	}

	fmt.Println("✅ Daemon is running")
	return nil
}