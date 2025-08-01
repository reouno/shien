package commands

import "shien/internal/rpc"

// Command represents a CLI command
type Command interface {
	Name() string
	Description() string
	Usage() string
	Execute(client *rpc.Client, args []string) error
}

// Registry manages available commands
type Registry struct {
	commands map[string]Command
}

// NewRegistry creates a new command registry
func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]Command),
	}
}

// Register adds a new command to the registry
func (r *Registry) Register(cmd Command) {
	r.commands[cmd.Name()] = cmd
}

// Get returns a command by name
func (r *Registry) Get(name string) (Command, bool) {
	cmd, exists := r.commands[name]
	return cmd, exists
}

// List returns all registered commands
func (r *Registry) List() map[string]Command {
	return r.commands
}