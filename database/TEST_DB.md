# Test Database Setup

This directory contains scripts and configuration for setting up a local test database.

## Quick Start

### Option 1: Using Docker (Recommended)

```bash
# Start PostgreSQL in Docker
docker compose -f docker-compose.test.yml up -d

# Wait for PostgreSQL to be ready (5-10 seconds)
sleep 5

# Setup test database
make test-db

# Run tests
make test

# Stop PostgreSQL when done
docker compose -f docker-compose.test.yml down
```

### Option 2: Using Local PostgreSQL

If you have PostgreSQL installed locally:

```bash
# Setup test database
make test-db

# Run tests
make test
```

## What `make test-db` Does

1. Checks if PostgreSQL is running
2. Drops existing `kasir_test` database (if any)
3. Creates fresh `kasir_test` database
4. Runs migrations on the test database

## Test Database Configuration

The tests expect:
- **Host:** localhost
- **Port:** 5432
- **User:** postgres
- **Password:** postgres
- **Database:** kasir_test
- **SSL Mode:** disable

## Troubleshooting

### PostgreSQL not running

**macOS (Homebrew):**
```bash
brew services start postgresql@14
```

**Linux:**
```bash
sudo systemctl start postgresql
```

**Docker:**
```bash
docker compose -f docker-compose.test.yml up -d
```

### Connection refused

Make sure PostgreSQL is listening on port 5432:
```bash
# Check if port is in use
lsof -i :5432

# Or use Docker
docker compose -f docker-compose.test.yml ps
```

### Permission denied

Update the password in `setup_test_db.sh` to match your PostgreSQL setup.

## CI/CD Integration

For CI/CD pipelines, use the Docker Compose setup:

```yaml
# GitHub Actions example
services:
  postgres:
    image: postgres:14-alpine
    env:
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: kasir_test
    ports:
      - 5432:5432
    options: >-
      --health-cmd pg_isready
      --health-interval 10s
      --health-timeout 5s
      --health-retries 5
```

Then run:
```bash
make test-db
make test
```
