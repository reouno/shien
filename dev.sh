#!/bin/bash

# Development environment setup script for Shien
# This script runs shien with an isolated development data directory

# Set the project root directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Development data directory
DEV_DATA_DIR="${PROJECT_ROOT}/.dev/data"

# Create development data directory
mkdir -p "${DEV_DATA_DIR}"

echo "Development environment:"
echo "  Data directory: ${DEV_DATA_DIR}"
echo ""

# Function to add --data-dir flag to commands
run_with_data_dir() {
    local cmd=$1
    shift
    
    case "$cmd" in
        ./shien|shien)
            # For daemon, add --data-dir flag
            "$cmd" --data-dir "${DEV_DATA_DIR}" "$@"
            ;;
        shienctl)
            # Use local build if available, otherwise system shienctl
            if [ -x "./shienctl" ]; then
                ./shienctl --data-dir "${DEV_DATA_DIR}" "$@"
            else
                echo "Error: ./shienctl not found. Run './dev.sh make build-all' first."
                exit 1
            fi
            ;;
        ./shienctl)
            # For local shienctl, add --data-dir flag
            "$cmd" --data-dir "${DEV_DATA_DIR}" "$@"
            ;;
        make)
            # For make commands, check if it's 'make run'
            if [ "$1" = "run" ]; then
                # Build and run with data-dir flag
                make build && ./shien --data-dir "${DEV_DATA_DIR}"
            else
                # Other make commands run normally
                "$cmd" "$@"
            fi
            ;;
        *)
            # Other commands run normally
            "$cmd" "$@"
            ;;
    esac
}

# Run the command passed as arguments
if [ $# -gt 0 ]; then
    run_with_data_dir "$@"
else
    echo "Usage: ./dev.sh <command> [args...]"
    echo ""
    echo "Examples:"
    echo "  ./dev.sh make run              # Build and run daemon in dev mode"
    echo "  ./dev.sh ./shien               # Run the daemon with dev data dir"
    echo "  ./dev.sh shienctl status       # Run shienctl with dev data dir"
    echo "  ./dev.sh shienctl activity -today"
    echo ""
    echo "The development environment uses: ${DEV_DATA_DIR}"
    echo "Production uses: ~/.config/shien"
fi