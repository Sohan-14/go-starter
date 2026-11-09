# ==========================================
# PHONY TARGETS
# ==========================================

.PHONY: all help run build test check fmt vet tidy clean

# ==========================================
# VARIABLES
# ==========================================

# Colors for help output
BLUE := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RESET := \033[0m

APP_NAME := go-starter
API_PATH := ./cmd/api
BIN_API_PATH := ./bin/api
BIN_DIR := ./bin


# Check if .env file exists and export its variables for local commands
ifneq (,$(wildcard .env))
    include .env
    export
endif

# ==========================================
# DEVELOPER COMMANDS
# ==========================================

# Show available commands
help:
	@echo "$(GREEN) Available commands for $(APP_NAME):$(RESET)"
	@echo "$(BLUE)  make run      - Run the application locally$(RESET)"
	@echo "$(BLUE)  make build    - Build the binary into ./bin$(RESET)"
	@echo "$(BLUE)  make test     - Run tests with race detector and coverage$(RESET)"
	@echo "$(BLUE)  make check    - Run all mandatory quality gates (fmt, vet, test)$(RESET)"
	@echo "$(BLUE)  make fmt      - Format the code$(RESET)"
	@echo "$(BLUE)  make tidy     - Clean up go.mod and go.sum$(RESET)"
	@echo "$(BLUE)  make clean    - Remove build artifacts and test cache$(RESET)"

# Run the application locally
run:
	@echo "$(GREEN)Running the api $(APP_NAME):$(RESET)"
	go run $(API_PATH)

# Build the binary
build:
	@echo "$(GREEN)Start building $(APP_NAME):$(RESET)"
	go build -o $(BIN_API_PATH) $(API_PATH)

# Run tests with race detector and coverage
test:
	@echo "$(YELLOW)Start testing $(APP_NAME):$(RESET)"
	go test -race -cover ./...

# ==========================================
# QUALITY GATES (CI & LOCAL)
# ==========================================

# Check formatting without modifying files (fails if unformatted)
fmt-check:
	@test -z $(shell gofmt -l .) || (echo "Code is not formatted. Run 'make fmt'"; exit 1)

# Run go vet (standard Go static analysis)
vet:
	go vet ./...

# The mandatory quality gate: format check, vet, and tests
# Note: Add 'golangci-lint run' to this list later if you adopt a heavier linter.
check: fmt-check vet test
	@echo "All quality gates passed."

# Format the code
fmt:
	gofmt -w .

# Tidy dependencies
tidy:
	go mod tidy

# ==========================================
# CLEANUP
# ==========================================

# Remove build artifacts and test cache
clean:
	rm -rf $(BIN_DIR)
	go clean -testcache
