all: build

build:
	@echo "Building..."
	go mod tidy
	go mod download
	@go build -o bin/$(shell basename $(PWD)) main.go

run:
	@echo "Running..."
	./bin/$(shell basename $(PWD))