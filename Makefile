# mochi - Development Environment Manager
# Makefile for development, testing, and CI workflows

.PHONY: help build install clean test test-unit test-integration test-e2e test-all lint fmt fmt-check docker-build docker-push release

# Default target
help: ## Show this help message
	@echo "mochi Development Makefile"
	@echo ""
	@echo "Development targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'
	@echo ""
	@echo "Testing targets:"
	@echo "  test-unit          Run unit tests with coverage"
	@echo "  test-integration   Run integration tests"
	@echo "  test-e2e           Run end-to-end tests"
	@echo "  test-all           Run all tests (unit + integration + e2e)"
	@echo ""
	@echo "CI targets:"
	@echo "  ci-check           Run all checks (lint, fmt, test)"
	@echo "  ci-build           Build for multiple platforms"
	@echo "  ci-docker          Build and push Docker image"

# Get version from git tag or use dev
VERSION := $(shell git describe --tags --always 2>/dev/null || echo "dev")
BUILDDATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -X main.version=$(VERSION) -X main.buildDate=$(BUILDDATE)

# Development
build: ## Build mochi binary
	go build -ldflags "$(LDFLAGS)" -o bin/mochi ./cmd/mochi

install: build ## Install mochi to ~/.local/bin
	cp bin/mochi ~/.local/bin/mochi
	chmod +x ~/.local/bin/mochi

ci-build: ## Build for multiple platforms
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/mochi-linux-amd64 ./cmd/mochi
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/mochi-linux-arm64 ./cmd/mochi
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/mochi-darwin-amd64 ./cmd/mochi
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/mochi-darwin-arm64 ./cmd/mochi
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/mochi-windows-amd64.exe ./cmd/mochi

ci-docker: docker-build docker-push ## Build and push Docker image

# Release
release: ci-check ci-build ## Create release artifacts
	@echo "Release artifacts created in dist/"

# Development helpers
dev-setup: ## Set up development environment
	go mod download
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

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
