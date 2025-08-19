all: build

build:
	@echo "Building..."
	@go mod tidy
	@go mod download
	@go build -o bin/$(shell basename $(PWD)) ./cmd

build_alone:
	@go build -o bin/$(shell basename $(PWD)) ./cmd

docker:
	@docker build -t ghcr.io/escape-ship/gatewaysrv:latest .

pushall:
	@docker build -t ghcr.io/escape-ship/gatewaysrv:latest .
	@docker push ghcr.io/escape-ship/gatewaysrv:latest

run:
	@echo "Running..."
	@./bin/$(shell basename $(PWD))

run_dev: ## Run with development configuration
	@echo "Running in development mode..."
	@export GATEWAY_AUTH_JWT_SECRET="your-secret-key-here" && \
	export GATEWAY_APP_LOG_LEVEL="debug" && \
	./bin/$(shell basename $(PWD))

test: ## Run tests
	@echo "Running tests..."
	@go test ./...

test_coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

validate_config: ## Validate configuration
	@echo "Validating configuration..."
	@export GATEWAY_AUTH_JWT_SECRET="test-secret" && \
	go run ./cmd -validate-config

help: ## Show help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: all init build build_alone pushall proto_gen tool_update tool_download run run_dev test test_coverage lint fmt clean validate_config help