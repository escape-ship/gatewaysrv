all: build

init:
	@echo "Initializing..."
	@$(MAKE) tool_download

build:
	@echo "Building..."
	@go mod tidy
	@go mod download
	@$(MAKE) proto_gen
	@go build -o bin/$(shell basename $(PWD)) ./cmd

build_alone:
	@go build -o bin/$(shell basename $(PWD)) ./cmd

proto_gen:
	@echo "Generating proto..."
	@cd proto && \
	buf dep update && \
	buf generate

tool_update:
	@echo "Updating tools..."
	@go get -modfile=tools.mod -tool github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	@go get -modfile=tools.mod -tool github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	@go get -modfile=tools.mod -tool google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go get -modfile=tools.mod -tool google.golang.org/protobuf/cmd/protoc-gen-go@latest

tool_download:
	@echo "Downloading tools..."
	@go install -modfile=tools.mod tool
	@go install github.com/bufbuild/buf/cmd/buf@latest

run:
	@echo "Running..."
	@./bin/$(shell basename $(PWD))

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci