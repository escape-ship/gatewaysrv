# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Build & Run
```bash
# Full build (includes proto generation)
make build

# Build without proto generation
make build_alone

# Run with development configuration
make run_dev

# Run built binary
make run
```

### Testing & Quality
```bash
# Run all tests
make test

# Run tests with coverage report
make test_coverage

# Run linter
make lint

# Format code
make fmt

# Validate configuration
make validate_config
```

### Protocol Buffers
```bash
# Generate protobuf code (included in build)
make proto_gen

# Update protobuf tools
make tool_update

# Download tools (run after cloning)
make tool_download
```

### Docker
```bash
# Build and push Docker image
make pushall
```

## Architecture Overview

This is a **gRPC Gateway service** that acts as an HTTP-to-gRPC proxy for microservices in the escape-ship ecosystem.

### Core Components

**Main Application Flow:**
- `cmd/main.go` → `internal/app/app.go` → `internal/gateway/gateway.go`
- Configuration loading with Viper from `config.yaml` and environment variables
- Graceful shutdown with signal handling (30s timeout)

**Gateway Pattern:**
- Uses grpc-gateway to translate HTTP REST calls to gRPC calls
- Registers handlers for 4 backend services: account, product, payment, order
- Each service has its own host/port configuration in `config.yaml`

**Middleware Chain (applied in reverse order):**
1. Recovery (outermost)
2. Auth (JWT validation)
3. CORS
4. Logging (innermost)

**Service Registry:**
- Account service: `accountsrv:8081`
- Product service: `productsrv:8082`
- Order service: `ordersrv:8083`
- Payment service: `paymentsrv:8084`

### Key Architecture Patterns

**Clean Architecture:**
- `cmd/` - Application entry point
- `internal/app/` - Application orchestration layer
- `internal/gateway/` - Gateway business logic
- `internal/middleware/` - HTTP middleware implementations
- `config/` - Configuration management
- `pkg/` - Shared utilities (logger, errors)

**Dependency Injection:**
- Config passed down through layers
- Logger injected into all components
- Gateway instance created and managed by App

**Error Handling:**
- Custom error types in `pkg/errors/`
- Structured logging with slog
- Graceful error propagation with context

## Configuration

### Environment Variables
Set these for development:
```bash
export GATEWAY_AUTH_JWT_SECRET="your-secret-key-here"
export GATEWAY_APP_LOG_LEVEL="debug"
```

### Service Discovery
Services are configured in `config.yaml` under the `services` section. Each service requires `host` and `port` values.

### Authentication
JWT secret is required and loaded from `GATEWAY_AUTH_JWT_SECRET` environment variable.

## Development Workflow

1. **Setup:** Run `make init` after cloning to download tools
2. **Development:** Use `make run_dev` for local development with debug logging
3. **Testing:** Run `make test_coverage` to ensure adequate test coverage
4. **Code Quality:** Run `make lint` and `make fmt` before committing
5. **Build:** Use `make build` for full build including proto generation

## Dependencies

**Core Framework:**
- grpc-gateway v2 for HTTP-to-gRPC translation
- Viper for configuration management
- slog for structured logging

**External Services:**
- Depends on `github.com/escape-ship/protos` for service definitions
- Connects to 4 backend gRPC services (account, product, order, payment)

**Development Tools:**
- golangci-lint for code quality
- buf for protocol buffer management
- Various protoc generators for gRPC code generation