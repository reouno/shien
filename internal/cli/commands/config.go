package commands

import (
	"encoding/json"
	"fmt"

	"shien/internal/rpc"
)

// ConfigCommand handles configuration display
type ConfigCommand struct{}

// NewConfigCommand creates a new config command
func NewConfigCommand() *ConfigCommand {
	return &ConfigCommand{}
}

// Name returns the command name
func (c *ConfigCommand) Name() string {
	return "config"
}

// Description returns the command description
func (c *ConfigCommand) Description() string {
	return "Show current configuration"
}

// Usage returns the command usage
func (c *ConfigCommand) Usage() string {
	return "config"
}

// Execute runs the config command
func (c *ConfigCommand) Execute(client *rpc.Client, args []string) error {
	resp, err := client.Call(rpc.MethodGetConfig, nil)
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("error: %s", resp.Error)
	}

	// Pretty print config
	data, err := json.MarshalIndent(resp.Data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format config: %w", err)
	}

	fmt.Println("Current Configuration")
	fmt.Println("====================")
	fmt.Println(string(data))

	return nil
}