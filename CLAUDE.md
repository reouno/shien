# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Shien (支援) is a Go-based daemon application designed to support knowledge workers. The application runs as a background service with graceful shutdown handling.

## Architecture

The codebase follows a standard Go project layout:
- `cmd/shien/main.go`: Entry point that initializes the daemon and system tray
- `internal/daemon/daemon.go`: Core daemon implementation with context-based cancellation
- `internal/tray/tray.go`: System tray UI implementation using getlantern/systray
- `internal/ui/display.go`: CLI display utilities for terminal output
- `internal/database/`: Database layer with repository pattern
  - `repository/`: Domain-specific repositories (activity, etc.)
  - `migrations/`: Database migrations organized by version
  - `models.go`: Data models
  - `interfaces.go`: Repository interfaces
- The daemon runs as a menu bar application with notification support
- Records activity every 5 minutes to track app usage

## Common Commands

### Build and Run
```bash
make build-all    # Build both daemon and CLI
make run          # Build and run the daemon
make install-all  # Install both to $GOPATH/bin
```

### Using the CLI
```bash
# Check if daemon is running
shienctl ping

# Show daemon status
shienctl status

# View activity logs
shienctl activity -today
shienctl activity -from 2024-01-01 -to 2024-01-31

# Show configuration
shienctl config
```

### Development
```bash
make test         # Run tests (currently no tests implemented)
make mod-tidy     # Clean up go.mod and go.sum
make clean        # Remove binary and clean cache
```

### Running the Application
```bash
./shien           # Run the built binary directly
```

## Development Notes

- The project uses Go 1.24.5
- Build flags include `-ldflags="-s -w"` for smaller binaries
- The daemon pattern uses context for cancellation and goroutines for concurrent operations
- Signal handling is implemented for graceful shutdown on SIGINT/SIGTERM