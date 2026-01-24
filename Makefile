.PHONY: help build run test coverage dev clean docs

help:
	@echo "Available targets:"
	@echo "  make build     - Build the application"
	@echo "  make run       - Run the application"
	@echo "  make test      - Run tests"
	@echo "  make coverage  - Run tests with coverage"
	@echo "  make dev       - Run with hot reload (requires air)"
	@echo "  make docs      - Generate OpenAPI documentation"
	@echo "  make clean     - Clean build artifacts"

build:
	go build -o kasir-api

run:
	go run .

test:
	go test -v -cover ./...

coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

dev:
	air

docs:
	~/tools/go/bin/swag init --v3.1 --outputTypes yaml,json

clean:
	rm -f kasir-api coverage.out coverage.html
