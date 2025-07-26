# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Shien (支援) is a Go-based daemon application designed to support knowledge workers. The application runs as a background service with graceful shutdown handling.

## Architecture

The codebase follows a standard Go project layout:
- `cmd/shien/main.go`: Entry point that initializes and manages the daemon lifecycle
- `internal/daemon/daemon.go`: Core daemon implementation with context-based cancellation
- The daemon runs a heartbeat every 30 seconds and responds to SIGINT/SIGTERM signals for graceful shutdown

## Common Commands

### Build and Run
```bash
make build        # Build the binary
make run          # Build and run the binary
make install      # Install to $GOPATH/bin
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