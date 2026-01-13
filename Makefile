# rightsizer Makefile

# Variables
BINARY_NAME := rightsizer
DOCKER_IMAGE := nextroll/rightsizer
GO_FILES := $(shell find . -name '*.go' -type f)

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	go build -o $(BINARY_NAME) .

# Generate mocks
.PHONY: generate
generate:
	go generate ./...

# Run tests
.PHONY: test
test: generate
	go test ./...

# Clean build artifacts
.PHONY: clean
clean:
	rm -f $(BINARY_NAME)

# Build Docker image
.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_IMAGE) .

# Run with Docker
.PHONY: docker-run
docker-run:
	docker run --rm $(DOCKER_IMAGE) $(BINARY_NAME) $(ARGS)

# Install dependencies
.PHONY: deps
deps:
	go mod download
	go mod tidy

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Vet code
.PHONY: vet
vet:
	go vet ./...

# Run linter (requires golangci-lint)
.PHONY: lint
lint:
	golangci-lint run

# Full check: format, vet, generate, test
.PHONY: check
check: fmt vet generate test

# Bump version
.PHONY: version-bump
version-bump:
	@if [ -z "$(VERSION)" ]; then \
		echo "Usage: make version-bump VERSION=x.y.z"; \
		exit 1; \
	fi
	@echo "Bumping version to $(VERSION)..."
	@sed -i '' 's/Version:.*"[0-9][0-9]*\.[0-9][0-9]*\.[0-9][0-9]*"/Version: "$(VERSION)"/' main.go
	@echo "Version updated to $(VERSION) in main.go"
	@# Extract major and minor versions for workflow tags
	@MAJOR=$$(echo $(VERSION) | cut -d. -f1); \
	MINOR=$$(echo $(VERSION) | cut -d. -f2); \
	sed -i '' "s/nextroll\/rightsizer:v[0-9][0-9]*/nextroll\/rightsizer:v$$MAJOR/" .github/workflows/ship-version.yml; \
	sed -i '' "s/nextroll\/rightsizer:v[0-9][0-9]*\.[0-9][0-9]*/nextroll\/rightsizer:v$$MAJOR.$$MINOR/" .github/workflows/ship-version.yml; \
	sed -i '' "s/nextroll\/rightsizer:v[0-9][0-9]*\.[0-9][0-9]*\.[0-9][0-9]*/nextroll\/rightsizer:v$(VERSION)/" .github/workflows/ship-version.yml
	@echo "Running go fmt..."
	go fmt ./...
	@echo "Version updated to $(VERSION) in .github/workflows/ship-version.yml"
	

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build        - Build the rightsizer binary"
	@echo "  generate     - Generate mocks using go generate"
	@echo "  test         - Run tests (generates mocks first)"
	@echo "  clean        - Remove build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run with Docker (use ARGS='cluster service')"
	@echo "  deps         - Download and tidy Go modules"
	@echo "  fmt          - Format Go code"
	@echo "  vet          - Run go vet"
	@echo "  lint         - Run golangci-lint"
	@echo "  check        - Run fmt, vet, generate, and test"
	@echo "  version-bump - Bump version (use VERSION=x.y.z)"
	@echo "  help         - Show this help message"
