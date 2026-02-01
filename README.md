# Bootcamp Jago Golang Dasar

## Kasir API - Clean Architecture Implementation

### Description
Kasir API is a REST API application for managing products and categories. Built using Go with **Clean Architecture**, **SOLID principles**, and **Test-Driven Development (TDD)**.

### Features
- ✅ CRUD (Create, Read, Update, Delete) Products
- ✅ CRUD Categories
- ✅ Dual storage: In-memory & PostgreSQL/Supabase
- ✅ Database migrations
- ✅ Configuration management (Koanf)
- ✅ Graceful shutdown
- ✅ Health check endpoint
- ✅ Comprehensive test coverage
- ✅ SOLID principles implementation
- ✅ Dependency Injection

### Architecture

```
cmd/
└── api/
    └── main.go              # Entry point with DI container
internal/
├── config/                  # Configuration management (Koanf)
│   ├── config.go
│   └── config_test.go
├── database/                # Database connection & migrations
│   ├── postgres.go
│   ├── migrate.go
│   └── *_test.go
├── model/                   # Domain entities with validation
│   ├── product.go
│   ├── category.go
│   ├── errors.go
│   └── *_test.go
├── dto/                     # Data Transfer Objects
│   └── product.go
├── repository/              # Data access layer
│   ├── interfaces.go        # Repository interfaces (ISP)
│   ├── memory/              # In-memory implementation
│   │   ├── product.go
│   │   ├── category.go
│   │   └── *_test.go
│   └── postgres/            # PostgreSQL implementation
│       ├── product.go
│       ├── category.go
│       └── *_test.go
├── service/                 # Business logic layer
│   ├── product.go
│   ├── category.go
│   └── *_test.go
└── handler/                 # HTTP handlers
    ├── product.go
    ├── category.go
    ├── health.go
    └── routes.go
pkg/
└── httputil/                # HTTP utilities
    └── response.go
database/
├── migrations/              # Database migrations (goose)
│   └── *.sql
├── seeds/                   # Seed data
│   └── seed.sql
└── rls/                     # Row Level Security policies
    ├── enable_rls.sql
    └── disable_rls.sql
```

### SOLID Principles Applied

- **S (Single Responsibility)**: Each layer has a single responsibility
- **O (Open/Closed)**: Easy to extend without modifying existing code
- **L (Liskov Substitution)**: In-memory and PostgreSQL repos are interchangeable
- **I (Interface Segregation)**: Small and focused interfaces (Reader/Writer)
- **D (Dependency Inversion)**: Depend on abstractions, not concretions

### Installation

```bash
# Clone repository
git clone <repository-url>
cd kasir-api

# Install dependencies
go mod download

# Build application
make build
```

### Configuration

Application uses environment variables with `APP_` prefix:

**Server Configuration:**
| Variable | Default | Description |
|----------|---------|-------------|
| `APP_SERVER_HOST` | `localhost` | Server host |
| `APP_SERVER_PORT` | `:8300` | Server port |
| `APP_SERVER_READTIMEOUT` | `10s` | Read timeout |
| `APP_SERVER_WRITETIMEOUT` | `10s` | Write timeout |

**Database Configuration (Optional):**
| Variable | Default | Description |
|----------|---------|-------------|
| `APP_DATABASE_HOST` | - | Database host |
| `APP_DATABASE_PORT` | `5432` | Database port |
| `APP_DATABASE_USER` | - | Database user |
| `APP_DATABASE_PASSWORD` | - | Database password |
| `APP_DATABASE_DBNAME` | - | Database name |
| `APP_DATABASE_SSLMODE` | `require` | SSL mode |
| `APP_DATABASE_MAXCONNS` | `25` | Max connections |
| `APP_DATABASE_MINCONNS` | `5` | Min connections |

**Note:** If database is not configured, application will use in-memory storage.

### Running the Application

#### Using Makefile (Recommended)
```bash
# View all commands
make help

# Build application
make build

# Run application
make run

# Run with PostgreSQL (set env vars first)
export APP_DATABASE_HOST=your-db-host.supabase.co
export APP_DATABASE_USER=postgres
export APP_DATABASE_PASSWORD=your-password
export APP_DATABASE_DBNAME=postgres
make run

# Run database migrations
make migrate

# Run tests
make test

# Run tests with coverage
make coverage

# Code quality check
make audit

# Development with hot reload (requires air)
make dev

# Clean build artifacts
make clean
```

#### Using Go
```bash
# Development (in-memory)
go run ./cmd/api/

# Build and run
go build -o bin/api ./cmd/api/
./bin/api

# Run migrations
./bin/api migrate

# With environment variables
APP_SERVER_PORT=:9000 go run ./cmd/api/
```

### Database Setup (PostgreSQL/Supabase)

1. **Create database** (if using local PostgreSQL):
```sql
CREATE DATABASE kasir_api;
```

2. **Configure environment variables**:
```bash
cp .env.example .env
# Edit .env with your database credentials
```

3. **Run migrations**:
```bash
make migrate
# or
./bin/api migrate
```

### API Documentation

Interactive API documentation is available via Scalar:

```bash
# Start the server
make run

# Open in browser
http://localhost:8300/docs
```

The documentation is generated from OpenAPI 3.0 specification and provides:
- Interactive API explorer
- Request/response examples
- Schema definitions
- Try-it-out functionality

### API Endpoints

#### Root
```bash
GET /
```
Response:
```json
{
  "message": "Kasir API"
}
```

#### Health Check
```bash
GET /health
```
Response:
```json
{
  "status": "OK",
  "message": "API is running"
}
```

#### Products

**Get All Products**
```bash
GET /api/products
```

**Get Product by ID**
```bash
GET /api/products/{id}
```

**Create Product**
```bash
POST /api/products
Content-Type: application/json

{
  "name": "Indomie Goreng",
  "price": 3500,
  "stock": 100
}
```

**Update Product**
```bash
PUT /api/products/{id}
Content-Type: application/json

{
  "name": "Indomie Goreng",
  "price": 4000,
  "stock": 150
}
```

**Delete Product**
```bash
DELETE /api/products/{id}
```

#### Categories

**Get All Categories**
```bash
GET /api/categories
```

**Get Category by ID**
```bash
GET /api/categories/{id}
```

**Create Category**
```bash
POST /api/categories
Content-Type: application/json

{
  "name": "Food",
  "description": "Food category"
}
```

**Update Category**
```bash
PUT /api/categories/{id}
Content-Type: application/json

{
  "name": "Food",
  "description": "Food and beverage category"
}
```

**Delete Category**
```bash
DELETE /api/categories/{id}
```

### Testing

```bash
# Run all tests
go test -v ./...

# Run specific test
go test -v -run TestProductRepository

# Run with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Development

#### Hot Reload with Air

Install air:
```bash
go install github.com/air-verse/air@latest
```

Run:
```bash
make dev
# or
air
```

### Technologies
- Go 1.25+
- Standard library (net/http, encoding/json)
- Testing with httptest

### License
MIT
