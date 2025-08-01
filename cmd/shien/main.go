package main

import (
	"fmt"
	"log"
	"os"

	"shien/internal/cli/commands"
	"shien/internal/paths"
	"shien/internal/rpc"
	"shien/internal/version"
)

func main() {
	// Parse global flags first
	dataDir := parseGlobalFlags()

	// Set custom data directory if provided
	if dataDir != "" {
		if err := paths.SetDataDir(dataDir); err != nil {
			log.Fatalf("Failed to set data directory: %v", err)
		}
	}

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	commandName := os.Args[1]

	// Handle special commands that don't need client
	switch commandName {
	case "help", "-h", "--help":
		printUsage()
		return
	case "version", "--version", "-v":
		fmt.Printf("shien version %s\n", version.GetFullVersion())
		return
	}

	// Create RPC client
	client, err := rpc.NewClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Initialize command registry
	registry := commands.NewRegistry()
	registerCommands(registry)

	// Execute command
	if cmd, exists := registry.Get(commandName); exists {
		if err := cmd.Execute(client, os.Args[2:]); err != nil {
			log.Fatalf("%s: %v", commandName, err)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", commandName)
		printUsage()
		os.Exit(1)
	}
}

func parseGlobalFlags() string {
	var dataDir string
	newArgs := []string{os.Args[0]} // Keep program name

	i := 1
	for i < len(os.Args) {
		if os.Args[i] == "--data-dir" && i+1 < len(os.Args) {
			dataDir = os.Args[i+1]
			i += 2 // Skip both --data-dir and its value
		} else {
			newArgs = append(newArgs, os.Args[i])
			i++
		}
	}
	os.Args = newArgs

	return dataDir
}

func registerCommands(registry *commands.Registry) {
	registry.Register(commands.NewStatusCommand())
	registry.Register(commands.NewActivityCommand())
	registry.Register(commands.NewConfigCommand())
	registry.Register(commands.NewPingCommand())
}

func printUsage() {
	fmt.Println("Usage: shien [--data-dir <path>] <command> [options]")
	fmt.Println()
	fmt.Println("Global Options:")
	fmt.Println("  --data-dir <path>   Use custom data directory")
	fmt.Println("  --version, -v       Show version information")
	fmt.Println()
	fmt.Println("Commands:")
	
	// Create temporary registry to list commands
	registry := commands.NewRegistry()
	registerCommands(registry)
	
	// Display each command with its description
	commandList := registry.List()
	for _, cmd := range []string{"status", "activity", "config", "ping"} { // Maintain order
		if command, exists := commandList[cmd]; exists {
			fmt.Printf("  %-20s %s\n", command.Name(), command.Description())
			if command.Usage() != command.Name() {
				// Print detailed usage if different from simple name
				usage := command.Usage()
				fmt.Printf("    %s\n", usage)
			}
		}
	}
	
	fmt.Println("  version             Show version information")
	fmt.Println("  help                Show this help message")
}