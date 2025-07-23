# Makefile for AI Project Tutorial API Server
# Provides build automation for Go-based API server with protobuf support

# Variables
BINARY_NAME := apiserver
OUTPUT_DIR := bin
MAIN_PATH := ./cmd/apiserver
PROTO_DIR := proto
GO_VERSION := 1.21

# Build flags
LDFLAGS := -w -s
BUILD_FLAGS := -ldflags="$(LDFLAGS)"

# Colors for output
RED := \033[31m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
RESET := \033[0m

# Default target
.DEFAULT_GOAL := help

## Build targets

.PHONY: build
build: ## Build the apiserver binary
	@echo "$(BLUE)Building $(BINARY_NAME)...$(RESET)"
	@mkdir -p $(OUTPUT_DIR)
	@go build $(BUILD_FLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)✅ Build complete: $(OUTPUT_DIR)/$(BINARY_NAME)$(RESET)"

.PHONY: build-debug
build-debug: ## Build the apiserver binary with debug symbols
	@echo "$(BLUE)Building $(BINARY_NAME) with debug symbols...$(RESET)"
	@mkdir -p $(OUTPUT_DIR)
	@go build -gcflags="all=-N -l" -o $(OUTPUT_DIR)/$(BINARY_NAME)-debug $(MAIN_PATH)
	@echo "$(GREEN)✅ Debug build complete: $(OUTPUT_DIR)/$(BINARY_NAME)-debug$(RESET)"

## Development targets

.PHONY: run
run: build ## Build and run the apiserver
	@echo "$(BLUE)Starting $(BINARY_NAME)...$(RESET)"
	@./$(OUTPUT_DIR)/$(BINARY_NAME)

.PHONY: test
test: ## Run all tests
	@echo "$(BLUE)Running tests...$(RESET)"
	@go test -v -race ./...
	@echo "$(GREEN)✅ All tests passed$(RESET)"

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	@echo "$(BLUE)Running tests with coverage...$(RESET)"
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✅ Coverage report generated: coverage.html$(RESET)"

.PHONY: bench
bench: ## Run benchmarks
	@echo "$(BLUE)Running benchmarks...$(RESET)"
	@go test -bench=. -benchmem ./...

## Code quality targets

.PHONY: fmt
fmt: ## Format Go code
	@echo "$(BLUE)Formatting code...$(RESET)"
	@go fmt ./...
	@echo "$(GREEN)✅ Code formatted$(RESET)"

.PHONY: vet
vet: ## Run go vet
	@echo "$(BLUE)Running go vet...$(RESET)"
	@go vet ./...
	@echo "$(GREEN)✅ Go vet passed$(RESET)"

.PHONY: lint
lint: install-golangci-lint ## Run golangci-lint
	@echo "$(BLUE)Running golangci-lint...$(RESET)"
	@golangci-lint run ./...
	@echo "$(GREEN)✅ Linting passed$(RESET)"

.PHONY: check
check: fmt vet test ## Run all code quality checks
	@echo "$(GREEN)✅ All checks passed$(RESET)"

## Dependency management

.PHONY: deps
deps: ## Download Go module dependencies
	@echo "$(BLUE)Downloading dependencies...$(RESET)"
	@go mod download
	@echo "$(GREEN)✅ Dependencies downloaded$(RESET)"

.PHONY: tidy
tidy: ## Tidy Go module dependencies
	@echo "$(BLUE)Tidying module dependencies...$(RESET)"
	@go mod tidy
	@echo "$(GREEN)✅ Module dependencies tidied$(RESET)"

.PHONY: verify
verify: ## Verify Go module dependencies
	@echo "$(BLUE)Verifying module dependencies...$(RESET)"
	@go mod verify
	@echo "$(GREEN)✅ Module dependencies verified$(RESET)"

## Protobuf targets (prepared for future use)

.PHONY: proto
proto: install-protobuf-tools ## Generate protobuf code
	@echo "$(BLUE)Generating protobuf code...$(RESET)"
	@if [ -n "$$(find $(PROTO_DIR) -name '*.proto' 2>/dev/null)" ]; then \
		buf generate; \
		echo "$(GREEN)✅ Protobuf code generated$(RESET)"; \
	else \
		echo "$(YELLOW)⚠️  No .proto files found in $(PROTO_DIR)$(RESET)"; \
	fi

.PHONY: proto-lint
proto-lint: install-buf ## Lint protobuf files
	@echo "$(BLUE)Linting protobuf files...$(RESET)"
	@if [ -n "$$(find $(PROTO_DIR) -name '*.proto' 2>/dev/null)" ]; then \
		buf lint; \
		echo "$(GREEN)✅ Protobuf linting passed$(RESET)"; \
	else \
		echo "$(YELLOW)⚠️  No .proto files found in $(PROTO_DIR)$(RESET)"; \
	fi

## Tool installation

.PHONY: install-tools
install-tools: install-golangci-lint install-buf install-protobuf-tools ## Install all development tools
	@echo "$(GREEN)✅ All development tools installed$(RESET)"

.PHONY: install-golangci-lint
install-golangci-lint: ## Install golangci-lint
	@echo "$(BLUE)Installing golangci-lint...$(RESET)"
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		echo "$(GREEN)✅ golangci-lint installed$(RESET)"; \
	else \
		echo "$(YELLOW)⚠️  golangci-lint already installed$(RESET)"; \
	fi

.PHONY: install-buf
install-buf: ## Install buf CLI for protobuf management
	@echo "$(BLUE)Installing buf CLI...$(RESET)"
	@if ! command -v buf >/dev/null 2>&1; then \
		go install github.com/bufbuild/buf/cmd/buf@latest; \
		echo "$(GREEN)✅ buf CLI installed$(RESET)"; \
	else \
		echo "$(YELLOW)⚠️  buf CLI already installed$(RESET)"; \
	fi

.PHONY: install-protobuf-tools
install-protobuf-tools: ## Install protobuf generation tools
	@echo "$(BLUE)Installing protobuf tools...$(RESET)"
	@if ! command -v protoc-gen-go >/dev/null 2>&1; then \
		go install google.golang.org/protobuf/cmd/protoc-gen-go@latest; \
		echo "$(GREEN)✅ protoc-gen-go installed$(RESET)"; \
	else \
		echo "$(YELLOW)⚠️  protoc-gen-go already installed$(RESET)"; \
	fi
	@if ! command -v protoc-gen-go-grpc >/dev/null 2>&1; then \
		go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest; \
		echo "$(GREEN)✅ protoc-gen-go-grpc installed$(RESET)"; \
	else \
		echo "$(YELLOW)⚠️  protoc-gen-go-grpc already installed$(RESET)"; \
	fi

## Utility targets

.PHONY: clean
clean: ## Clean build artifacts and temporary files
	@echo "$(BLUE)Cleaning build artifacts...$(RESET)"
	@rm -rf $(OUTPUT_DIR)
	@rm -f coverage.out coverage.html
	@echo "$(GREEN)✅ Clean complete$(RESET)"

.PHONY: version
version: ## Show Go version and module info
	@echo "$(BLUE)Go version:$(RESET)"
	@go version
	@echo "$(BLUE)Module info:$(RESET)"
	@go list -m

.PHONY: env
env: ## Show Go environment
	@echo "$(BLUE)Go environment:$(RESET)"
	@go env

.PHONY: help
help: ## Show this help message
	@echo "$(BLUE)AI Project Tutorial API Server - Build Automation$(RESET)"
	@echo ""
	@echo "$(YELLOW)Usage:$(RESET)"
	@echo "  make <target>"
	@echo ""
	@echo "$(YELLOW)Targets:$(RESET)"
	@awk 'BEGIN {FS = ":.*##"; printf ""} /^[a-zA-Z_-]+:.*?##/ { printf "  $(GREEN)%-20s$(RESET) %s\n", $$1, $$2 } /^##@/ { printf "\n$(BLUE)%s$(RESET)\n", substr($$0, 5) } ' $(MAKEFILE_LIST) 