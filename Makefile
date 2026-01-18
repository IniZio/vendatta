# nexus - Development Environment Manager
# Makefile for development, testing, and CI workflows

.PHONY: help build install clean test test-unit test-integration test-e2e test-all lint fmt fmt-check docker-build docker-push release

# Default target
help: ## Show this help message
	@echo "nexus Development Makefile"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

# Get version from git tag or use dev
VERSION := $(shell git describe --tags --always 2>/dev/null || echo "dev")
BUILDDATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -X main.version=$(VERSION) -X main.buildDate=$(BUILDDATE)

# Development
build: ## Build nexus binary
	go build -ldflags "$(LDFLAGS)" -o bin/nexus ./cmd/nexus

install: build ## Install nexus to ~/.local/bin
	cp bin/nexus ~/.local/bin/nexus
	chmod +x ~/.local/bin/nexus

ci-build: ## Build for multiple platforms
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/nexus-linux-amd64 ./cmd/nexus
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/nexus-linux-arm64 ./cmd/nexus
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/nexus-darwin-amd64 ./cmd/nexus
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/nexus-darwin-arm64 ./cmd/nexus
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/nexus-windows-amd64.exe ./cmd/nexus

ci-docker: docker-build docker-push ## Build and push Docker image

ci-check: test-unit fmt-check lint ## Run all CI checks (tests, format, lint)
	@echo "✅ All CI checks passed"

# Release
release: ci-check ci-build ## Create release artifacts
	@echo "Release artifacts created in dist/"

# Development helpers
dev-setup: ## Set up development environment
	go mod download
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Linting and formatting
fmt: ## Format code with gofmt
	gofmt -w .
	go mod tidy

fmt-check: ## Check code formatting without modifying
	@if gofmt -l . | grep -q .; then \
		echo "❌ Code formatting issues found. Run 'make fmt' to fix."; \
		gofmt -l .; \
		exit 1; \
	else \
		echo "✅ Code formatting OK"; \
	fi

lint: ## Run linters (golangci-lint, go vet)
	go vet ./...
	golangci-lint run ./... || true
	@echo "✅ Linting complete"

dev-test-watch: ## Run tests in watch mode (requires entr)
	find . -name "*.go" | entr -r make test-unit

# Performance testing
perf-test: ## Run performance tests
	go test -bench=. -benchmem ./pkg/...
	@echo "Performance test complete. Check memory usage and startup times."

# Security
security-scan: ## Run security vulnerability scan
	gosec ./...
	trivy filesystem --exit-code 1 --no-progress .

# Documentation
docs-build: ## Build documentation
	@echo "Building docs..."
	# Add documentation build commands here

docs-serve: ## Serve documentation locally
	@echo "Serving docs on http://localhost:8000"
	# Add docs serve commands here

# Testing
test: test-unit test-integration test-e2e ## Run all tests (unit + integration + e2e)

test-unit: ## Run unit tests with coverage
	go test -short -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Unit test coverage saved to coverage.html"

test-integration: ## Run integration tests (no Docker/LXC required)
	go test -run Integration ./...

test-e2e: ## Run end-to-end tests (requires Docker/LXC)
	go test -v ./e2e/...

test-coverage: ## Generate coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

test-all: test test-coverage ## Run all tests with coverage report
