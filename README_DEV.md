# Development Environment Setup

This document explains how to run Shien in development mode without affecting your production installation.

## Quick Start

Use the `dev.sh` script to run Shien with an isolated data directory:

```bash
# Build binaries
./dev.sh make build      # Build daemon only
./dev.sh make build-all  # Build both daemon and CLI

# Run daemon in development mode
./dev.sh make run

# Or if already built
./dev.sh ./shien

# Use shienctl with development data
./dev.sh shienctl status
./dev.sh shienctl activity -today
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
./shien --data-dir /path/to/dev/data

# Use shienctl with same data directory
shienctl --data-dir /path/to/dev/data status
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
- `shien.sock` - Unix socket for IPC

## Note

The `.dev/` directory is already in `.gitignore`, so your development data won't be committed to the repository.