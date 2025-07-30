.PHONY: build build-all build-daemon build-cli run clean install install-all test dev-build-all

DAEMON_NAME=shien
CLI_NAME=shienctl
GO=go
GOFLAGS=-ldflags="-s -w"
DEV_TAGS=-tags dev

build: build-daemon

build-all: build-daemon build-cli

# Development build targets
dev-build-all:
	$(GO) build $(GOFLAGS) $(DEV_TAGS) -o $(DAEMON_NAME) ./cmd/shien
	$(GO) build $(GOFLAGS) $(DEV_TAGS) -o $(CLI_NAME) ./cmd/shienctl
	@echo "Built development version - data directory: .dev/data"

build-daemon:
	$(GO) build $(GOFLAGS) -o $(DAEMON_NAME) ./cmd/shien

build-cli:
	$(GO) build $(GOFLAGS) -o $(CLI_NAME) ./cmd/shienctl

run: build-daemon
	./$(DAEMON_NAME)

clean:
	rm -f $(DAEMON_NAME) $(CLI_NAME)
	$(GO) clean

install: install-all

install-all: build-all
	$(GO) install ./cmd/shien
	$(GO) install ./cmd/shienctl

install-daemon: build-daemon
	$(GO) install ./cmd/shien

install-cli: build-cli
	$(GO) install ./cmd/shienctl

test:
	$(GO) test -v ./...

mod-tidy:
	$(GO) mod tidy

help:
	@echo "Available targets:"
	@echo "  build        - Build the daemon (default)"
	@echo "  build-all    - Build both daemon and CLI (production)"
	@echo "  dev-build-all - Build both daemon and CLI (development mode)"
	@echo "  build-daemon - Build only the daemon"
	@echo "  build-cli    - Build only the CLI"
	@echo "  run          - Build and run the daemon"
	@echo "  clean        - Remove binaries and clean cache"
	@echo "  install      - Install both daemon and CLI to GOPATH/bin"
	@echo "  install-all  - Install both daemon and CLI"
	@echo "  install-daemon - Install only the daemon"
	@echo "  install-cli  - Install only the CLI"
	@echo "  test         - Run tests"
	@echo "  mod-tidy     - Clean up go.mod and go.sum"
	@echo "  help         - Show this help message"