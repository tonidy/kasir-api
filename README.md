# Bootcamp Jago Golang Dasar

## Task Session 1

### Deskripsi
Kasir API adalah aplikasi REST API sederhana untuk mengelola produk dan kategori. Dibangun menggunakan Go dengan arsitektur yang bersih dan terstruktur.

### Fitur
- CRUD (Create, Read, Update, Delete) Produk
- CRUD Kategori
- In-memory storage
- Graceful shutdown
- Health check endpoint
- Environment variable configuration

### Struktur Proyek
```
.
├── main.go           # Entry point aplikasi
├── server.go         # Server lifecycle management
├── route.go          # Route definitions
├── handler.go        # HTTP handlers
├── repository.go     # Data access layer
├── model.go          # Data models
├── constants.go      # Constants dan konfigurasi
├── *_test.go         # Unit tests
└── Makefile          # Build automation
```

### Instalasi

```bash
# Clone repository
git clone <repository-url>
cd kasir-api

# Install dependencies
go mod download
```

### Konfigurasi

Aplikasi menggunakan environment variables:

| Variable | Default | Deskripsi |
|----------|---------|-----------|
| `SERVER_HOST` | `localhost` | Host server |
| `SERVER_PORT` | `:8300` | Port server |

### Menjalankan Aplikasi

#### Menggunakan Go
```bash
# Development
go run .

# Build dan run
go build -o kasir-api
./kasir-api

# Dengan environment variables
SERVER_HOST=0.0.0.0 SERVER_PORT=:9000 go run .
```

#### Menggunakan Makefile
```bash
# Lihat semua perintah
make help

# Build aplikasi
make build

# Run aplikasi
make run

# Run tests
make test

# Run tests dengan coverage
make coverage

# Development dengan hot reload (requires air)
make dev

# Clean build artifacts
make clean
```

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
  "pesan": "API berjalan dengan baik"
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
  "name": "Makanan",
  "description": "Kategori makanan"
}
```

**Update Category**
```bash
PUT /api/categories/{id}
Content-Type: application/json

{
  "name": "Makanan",
  "description": "Kategori makanan dan minuman"
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

#### Hot Reload dengan Air

Install air:
```bash
go install github.com/air-verse/air@latest
```

Jalankan:
```bash
make dev
# atau
air
```

### Teknologi
- Go 1.25+
- Standard library (net/http, encoding/json)
- Testing dengan httptest

### License
MIT
