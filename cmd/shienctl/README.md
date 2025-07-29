# shienctl - Shien CLI Client

Command-line interface for interacting with the Shien daemon.

## Installation

```bash
make install-cli
# or
go install ./cmd/shienctl
```

## Usage

### Check daemon status
```bash
shienctl ping
shienctl status
```

### View activity logs
```bash
# Today's activity
shienctl activity -today

# Specific date range
shienctl activity -from 2024-01-01 -to 2024-01-31

# Last 24 hours (default)
shienctl activity
```

### Configuration
```bash
# View current configuration
shienctl config

# Update configuration (future feature)
# shienctl config set notification_enabled=true
```

## Architecture

The CLI communicates with the daemon via Unix socket located at:
- `~/.config/shien/shien.sock`

This ensures:
- Only the user who started the daemon can access it
- No network exposure
- Fast local communication
- Data consistency (daemon is the single source of truth)