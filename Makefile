.PHONY: build test clean help lint fmt run

# Build the lingo compiler
build:
	@echo "Building lingo compiler..."
	go build -o bin/lingo ./cmd/lingo/main.go
	go build -o bin/lingoctl ./cmd/lingoctl/main.go
	@echo "Build complete.  Binaries in ./bin/"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -cover ./...

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Run the compiler on an example file
run-example:
	@echo "Compiling example..."
	./bin/lingo -file examples/basic.lingo -out examples/basic.go -v

# Run lingoctl
ctl-lex:
	@echo "Lexing example..."
	./bin/lingoctl -cmd lex -file examples/basic.lingo

ctl-parse:
	@echo "Parsing example..."
	./bin/lingoctl -cmd parse -file examples/basic.lingo

# Build and run tests
all: clean deps build test

# Help
help:
	@echo "Lingo - TypeScript-like Meta-Language for Go"
	@echo ""
	@echo "Available targets:"
	@echo "  make build          - Build the lingo compiler and lingoctl"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make fmt            - Format code with gofmt"
	@echo "  make lint           - Lint code with golangci-lint"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make deps           - Download and tidy dependencies"
	@echo "  make run-example    - Compile an example file"
	@echo "  make ctl-lex        - Lex an example file"
	@echo "  make ctl-parse      - Parse an example file"
	@echo "  make all            - Clean, install deps, build, and test"
	@echo "  make help           - Show this help message"
