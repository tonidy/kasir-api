.PHONY: help build run test coverage dev clean docs audit migrate seed rls-on rls-off

help:
	@echo "Available targets:"
	@echo "  make build     - Build the application"
	@echo "  make run       - Run the application"
	@echo "  make test      - Run tests"
	@echo "  make coverage  - Run tests with coverage"
	@echo "  make dev       - Run with hot reload (requires air)"
	@echo "  make docs      - Generate OpenAPI documentation"
	@echo "  make audit     - Tidy, format, vet, and run static check"
	@echo "  make migrate   - Run database migrations"
	@echo "  make seed      - Seed database with sample data"
	@echo "  make rls-on    - Enable Row Level Security"
	@echo "  make rls-off   - Disable Row Level Security"
	@echo "  make clean     - Clean build artifacts"

build:
	go build -o bin/kasir-api ./cmd/api/

run:
	go run ./cmd/api/

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

migrate:
	go run ./cmd/api migrate

seed:
	go run ./cmd/api seed

rls-on:
	go run ./cmd/api rls on

rls-off:
	go run ./cmd/api rls off

audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	@echo 'Running static check...'
	staticcheck ./...

clean:
	rm -rf bin/ coverage.out coverage.html
