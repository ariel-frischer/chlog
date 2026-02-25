.PHONY: help install i test lint format clean build run go-install

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Download dependencies
	go mod download

i: install ## Alias for install

go-install: ## Install chlog to GOPATH/bin
	go install ./cmd/chlog/

test: ## Run tests
	go test ./...

test-v: ## Run tests (verbose)
	go test -v ./...

lint: ## Run linters
	golangci-lint run

format: ## Format code
	go fmt ./...

clean: ## Clean build artifacts
	go clean
	rm -rf bin/

build: ## Build binary
	go build -o bin/chlog ./cmd/chlog/

run: ## Run main package
	go run ./cmd/chlog/

##@ Changelog
changelog-sync: build ## Regenerate CHANGELOG.md from CHANGELOG.yaml
	@./bin/chlog sync

changelog-check: build ## Validate CHANGELOG.md matches CHANGELOG.yaml
	@./bin/chlog check
