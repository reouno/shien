# Development Environment Setup

This document explains how to run Shien in development mode without affecting your production installation.

## Quick Start

Use the development build commands to run Shien with an isolated data directory:

```bash
# Build binaries in development mode
make dev-build-all  # Build both daemon and CLI with dev tags

# Run daemon in development mode
./shien-service --data-dir .dev/data

# Use CLI with development data
./shien --data-dir .dev/data status
./shien --data-dir .dev/data activity -today
```

## How It Works

The development setup uses a separate data directory (`.dev/data/`) instead of the default `~/.config/shien/`. This allows you to:

1. Run development and production versions simultaneously
2. Test with clean data without affecting your real data
3. Experiment freely without risk

## Manual Usage

You can also manually specify the data directory:

```bash
# Run daemon with custom data directory
./shien-service --data-dir /path/to/dev/data

# Use CLI with same data directory
./shien --data-dir /path/to/dev/data status
```

## Priority Order

The data directory is determined in this order:
1. Command-line flag `--data-dir`
2. Environment variable `SHIEN_DATA_DIR`
3. Default `~/.config/shien`

## Files Created

When using a custom data directory, all files are created there:
- `config.json` - Configuration file
- `shien.db` - SQLite database
- `shien-service.sock` - Unix socket for IPC

## Note

The `.dev/` directory is already in `.gitignore`, so your development data won't be committed to the repository.