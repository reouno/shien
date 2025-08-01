# Shien CLI Architecture

This directory contains the modular CLI implementation for the shien command-line tool.

## Directory Structure

```
internal/cli/
├── commands/        # Command implementations
│   ├── command.go   # Command interface and registry
│   ├── activity.go  # Activity log command
│   ├── config.go    # Configuration display command
│   ├── ping.go      # Daemon health check command
│   └── status.go    # Daemon status command
└── display/         # Display utilities
    └── activity.go  # Activity report formatting
```

## Adding a New Command

To add a new command to shien:

1. **Create a new command file** in `internal/cli/commands/`:
   ```go
   // internal/cli/commands/mycommand.go
   package commands
   
   type MyCommand struct{}
   
   func NewMyCommand() *MyCommand {
       return &MyCommand{}
   }
   
   func (c *MyCommand) Name() string { return "mycommand" }
   func (c *MyCommand) Description() string { return "Does something useful" }
   func (c *MyCommand) Usage() string { return "mycommand [options]" }
   func (c *MyCommand) Execute(client *rpc.Client, args []string) error {
       // Implementation here
       return nil
   }
   ```

2. **Register the command** in `cmd/shien/main.go`:
   ```go
   func registerCommands(registry *commands.Registry) {
       // ... existing commands ...
       registry.Register(commands.NewMyCommand())
   }
   ```

3. **Add RPC method** if needed in `internal/rpc/methods.go`

See `template.go.example` for a complete example.

## Design Principles

- **Simplicity**: Keep the architecture simple and easy to understand
- **Modularity**: Each command is self-contained in its own file
- **Testability**: Commands implement a common interface for easy testing
- **Extensibility**: New commands can be added without modifying existing code

## Command Interface

All commands must implement the `Command` interface:

```go
type Command interface {
    Name() string                                    // Command name
    Description() string                             // Short description
    Usage() string                                   // Usage information
    Execute(client *rpc.Client, args []string) error // Command logic
}
```

## Display Utilities

The `display` package contains reusable formatting utilities. For example, `ActivityReporter` handles the formatting of activity logs with hourly breakdowns and visual bars.

When adding complex display logic, consider creating a dedicated formatter in the `display` package to keep commands focused on business logic.