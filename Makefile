.PHONY: help build run test coverage dev clean

help:
	@echo "Available targets:"
	@echo "  make build     - Build the application"
	@echo "  make run       - Run the application"
	@echo "  make test      - Run tests"
	@echo "  make coverage  - Run tests with coverage"
	@echo "  make dev       - Run with hot reload (requires air)"
	@echo "  make clean     - Clean build artifacts"

build:
	go build -o kasir-api

run:
	go run .

test:
	go test -v ./...

coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

dev:
	air

clean:
	rm -f kasir-api coverage.out coverage.html
