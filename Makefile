.PHONY: help install i test test-v test-coverage lint lint-go format clean build run go-install release patch minor major prep-release

MODULE_PATH=github.com/ariel-frischer/chlog
VERSION?=$(shell git tag --sort=-v:refname 2>/dev/null | head -1)
ifeq ($(VERSION),)
  VERSION=dev
endif
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS=-ldflags="-X ${MODULE_PATH}/internal/version.Version=${VERSION} \
                   -X ${MODULE_PATH}/internal/version.Commit=${COMMIT} \
                   -X ${MODULE_PATH}/internal/version.BuildDate=${BUILD_DATE} \
                   -s -w"

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Download dependencies
	go mod download

i: install ## Alias for install

go-install: ## Install chlog to GOPATH/bin
	go install ${LDFLAGS} ./cmd/chlog/

test: ## Run tests
	go test ./...

test-v: ## Run tests (verbose)
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -race -coverprofile=coverage.out ./...

lint: lint-go ## Run all linters

lint-go: ## Run Go linters
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, running go vet"; \
		go vet ./...; \
	fi

format: ## Format code
	go fmt ./...

clean: ## Clean build artifacts
	go clean
	rm -rf bin/ coverage.out

build: ## Build binary with version info
	go build ${LDFLAGS} -o bin/chlog ./cmd/chlog/

run: ## Run main package
	go run ${LDFLAGS} ./cmd/chlog/

##@ Release
prep-release: build ## Full release flow: checks, stamp changelog, sync, commit, tag, push (usage: make prep-release VERSION=v0.1.0)
	@if [ -z "$(VERSION)" ] || [ "$(VERSION)" = "dev" ]; then \
		echo "Error: VERSION is required (e.g., make prep-release VERSION=v0.1.0)"; \
		exit 1; \
	fi
	@echo "==> Pre-flight checks..."
	@test -z "$$(git status --porcelain)" || { echo "Error: working tree is dirty"; exit 1; }
	@echo "==> Running tests..."
	@$(MAKE) test
	@echo "==> Running lint..."
	@$(MAKE) lint
	@echo "==> Checking unreleased entries..."
	@./bin/chlog show unreleased | grep -q . || { echo "Error: no unreleased entries in CHANGELOG.yaml"; exit 1; }
	@echo "==> Stamping changelog: $(VERSION)..."
	@VER=$$(echo "$(VERSION)" | sed 's/^v//'); ./bin/chlog release "$$VER"
	@echo "==> Syncing CHANGELOG.md..."
	@./bin/chlog sync
	@echo "==> Committing changelog..."
	@git add CHANGELOG.yaml CHANGELOG.md
	@git commit -m "release: $(VERSION)"
	@echo "==> Tagging $(VERSION)..."
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@echo "==> Pushing to gh..."
	@git push gh main
	@git push gh $(VERSION)
	@echo "Done! Run 'goreleaser release --clean' for GitHub release with binaries."

release: ## Create a release tag and push (usage: make release VERSION=v1.0.0)
	@if [ "$(VERSION)" = "dev" ] || [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required (e.g., make release VERSION=v1.0.0)"; \
		exit 1; \
	fi
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push gh $(VERSION)

patch: ## Bump patch version and release
	$(eval CURRENT=$(shell git tag --sort=-v:refname | head -1 | sed 's/^v//'))
	$(eval NEXT=v$(shell echo $(CURRENT) | awk -F. '{printf "%d.%d.%d", $$1, $$2, $$3+1}'))
	@echo "Bumping $(CURRENT) -> $(NEXT)"
	$(MAKE) release VERSION=$(NEXT)

minor: ## Bump minor version and release
	$(eval CURRENT=$(shell git tag --sort=-v:refname | head -1 | sed 's/^v//'))
	$(eval NEXT=v$(shell echo $(CURRENT) | awk -F. '{printf "%d.%d.0", $$1, $$2+1}'))
	@echo "Bumping $(CURRENT) -> $(NEXT)"
	$(MAKE) release VERSION=$(NEXT)

major: ## Bump major version and release
	$(eval CURRENT=$(shell git tag --sort=-v:refname | head -1 | sed 's/^v//'))
	$(eval NEXT=v$(shell echo $(CURRENT) | awk -F. '{printf "%d.0.0", $$1+1}'))
	@echo "Bumping $(CURRENT) -> $(NEXT)"
	$(MAKE) release VERSION=$(NEXT)

##@ Changelog
changelog-sync: build ## Regenerate CHANGELOG.md from CHANGELOG.yaml
	@./bin/chlog sync

changelog-check: build ## Validate CHANGELOG.md matches CHANGELOG.yaml
	@./bin/chlog check
