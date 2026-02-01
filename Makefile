.PHONY: help build run test coverage dev clean docs audit migrate

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
