# Shien (支援)

The tool to support knowledge workers worldwide.

## Installation

### Homebrew (macOS)

```bash
brew tap reouno/shien
brew install shien
brew services start shien
```

### Verify Installation

```bash
shienctl ping
```

## Build and Run

### Build
```bash
make build
```

### Run
```bash
# Build and run
make run

# Or run the built binary directly
./shien
```

### Install
```bash
# Install to $GOPATH/bin
make install
```

### Other Commands
```bash
# Show help
make help

# Clean up
make clean
```
