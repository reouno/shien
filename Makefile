.PHONY: build run clean install test

BINARY_NAME=shien
GO=go
GOFLAGS=-ldflags="-s -w"

build:
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) ./cmd/shien

run: build
	./$(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)
	$(GO) clean

install: build
	$(GO) install ./cmd/shien

test:
	$(GO) test -v ./...

mod-tidy:
	$(GO) mod tidy

help:
	@echo "Available targets:"
	@echo "  build      - Build the binary"
	@echo "  run        - Build and run the binary"
	@echo "  clean      - Remove binary and clean cache"
	@echo "  install    - Install the binary to GOPATH/bin"
	@echo "  test       - Run tests"
	@echo "  mod-tidy   - Clean up go.mod and go.sum"
	@echo "  help       - Show this help message"