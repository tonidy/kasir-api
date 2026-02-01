#!/bin/bash

# Setup test database for local PostgreSQL
# Usage: ./database/setup_test_db.sh

set -e

DB_NAME="kasir_test"
DB_USER="postgres"
DB_PASSWORD="postgres"
DB_HOST="localhost"
DB_PORT="5432"

echo "Setting up test database: $DB_NAME"

# Check if PostgreSQL is running
if ! pg_isready -h $DB_HOST -p $DB_PORT -U $DB_USER > /dev/null 2>&1; then
    echo "❌ PostgreSQL is not running on $DB_HOST:$DB_PORT"
    echo ""
    echo "To start PostgreSQL:"
    echo "  macOS (Homebrew): brew services start postgresql@14"
    echo "  Linux: sudo systemctl start postgresql"
    echo "  Docker: docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=postgres postgres:14"
    exit 1
fi

# Drop database if exists
echo "Dropping existing database (if any)..."
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -c "DROP DATABASE IF EXISTS $DB_NAME;" postgres 2>/dev/null || true

# Create database
echo "Creating database..."
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -c "CREATE DATABASE $DB_NAME;" postgres

# Run migrations on test database
echo "Running migrations..."
export APP_DATABASE_HOST=$DB_HOST
export APP_DATABASE_PORT=$DB_PORT
export APP_DATABASE_USER=$DB_USER
export APP_DATABASE_PASSWORD=$DB_PASSWORD
export APP_DATABASE_DBNAME=$DB_NAME
export APP_DATABASE_SSLMODE=disable

go run ./cmd/api migrate

echo "✅ Test database setup complete!"
echo ""
echo "Connection details:"
echo "  Host: $DB_HOST"
echo "  Port: $DB_PORT"
echo "  Database: $DB_NAME"
echo "  User: $DB_USER"
echo ""
echo "Run tests with: make test"
