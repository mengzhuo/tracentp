.PHONY: build test clean install lint format help

# Default target
all: build

# Build the application
build:
	go build -o tracentp main.go

# Install the application
install:
	go install .

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run benchmarks
bench:
	go test -bench=. ./...

# Clean build artifacts
clean:
	rm -f tracentp
	rm -f coverage.out coverage.html
	rm -rf dist/

# Format code
format:
	go fmt ./...
	gofmt -s -w .

# Lint code
lint:
	golangci-lint run

# Run GoReleaser locally (for testing)
release-dry-run:
	goreleaser release --snapshot --clean --skip-publish

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o tracentp-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build -o tracentp-linux-arm64 main.go
	GOOS=linux GOARCH=riscv64 go build -o tracentp-linux-riscv64 main.go
	GOOS=freebsd GOARCH=amd64 go build -o tracentp-freebsd-amd64 main.go
	GOOS=freebsd GOARCH=arm64 go build -o tracentp-freebsd-arm64 main.go
	GOOS=darwin GOARCH=amd64 go build -o tracentp-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o tracentp-darwin-arm64 main.go
	GOOS=windows GOARCH=amd64 go build -o tracentp-windows-amd64.exe main.go

# Show help
help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  install      - Install the application"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  bench        - Run benchmarks"
	@echo "  clean        - Clean build artifacts"
	@echo "  format       - Format code"
	@echo "  lint         - Lint code"
	@echo "  release-dry-run - Test GoReleaser configuration"
	@echo "  build-all    - Build for multiple platforms"
	@echo "  help         - Show this help" 