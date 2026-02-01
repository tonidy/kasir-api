# Binary Size History

This document tracks the evolution of the application's binary size and dependencies over time.

## Size Comparison Table

| Phase | Date       | Binary Size | Change       | Dependencies | Key Features                   |
| ----- | ---------- | ----------- | ------------ | ------------ | ------------------------------ |
| 1     | 2025-01-25 | 7.9M        | -            | 0            | In-memory, stdlib only         |
| 2     | 2026-02-01 | 15M         | +7.1M (+90%) | 23           | PostgreSQL, migrations, config |

## Timeline

### Phase 1: Initial Implementation (Commit: initial)
**Date:** 2025-01-25 
**Branch:** main

**Binary Size:** 7.9M

**Dependencies:**
```go
module kasir-api

go 1.25.6
```

**Features:**
- In-memory storage only
- Standard library HTTP server
- Zero external dependencies
- Monolithic structure

---

### Phase 2: Clean Architecture + PostgreSQL (Commit: 0c833ed)
**Date:** 2026-02-01  
**Branch:** task-session-2

**Binary Size:** 15M (+7.1M, +90%)

**Dependencies Added:**
```go
require (
	github.com/jackc/pgx/v5 v5.8.0              // PostgreSQL driver
	github.com/knadh/koanf/parsers/dotenv v1.1.1 // .env parser
	github.com/knadh/koanf/providers/env/v2 v2.0.0 // Env provider
	github.com/knadh/koanf/providers/file v1.2.1  // File provider
	github.com/knadh/koanf/v2 v2.3.2             // Config management
	github.com/pressly/goose/v3 v3.26.0          // Database migrations
)
```

**Dependency Stats:**
- Direct: 6 packages
- Indirect: 17 packages
- Total: 23 packages
- go.mod: 32 lines
- go.sum: 77 lines

**Size Breakdown (estimated):**
- Base application: ~3M
- pgx/v5: ~4-5M
- koanf: ~1-2M
- goose: ~1-2M
- Other dependencies: ~2-3M

**Features Added:**
- ✅ Layered architecture (domain, repository, service, handler)
- ✅ PostgreSQL support with connection pooling
- ✅ In-memory fallback (dual storage)
- ✅ Configuration management (.env + environment variables)
- ✅ Database migrations (goose)
- ✅ Database seeding
- ✅ Row Level Security (RLS) support
- ✅ Dependency injection
- ✅ SOLID principles
- ✅ Comprehensive test coverage

---

## Optimization Options

If binary size becomes a concern, consider:

### 1. Build with Stripped Symbols
```bash
# Strip debug symbols and symbol table
go build -ldflags="-s -w" -o bin/kasir-api ./cmd/api/
```
**Expected size:** ~10M (33% reduction from 15M)

**What it does:**
- `-s`: Omit symbol table and debug info
- `-w`: Omit DWARF symbol table

**Verified on macOS:** 15M → 10M

### 2. Compress with UPX (Linux/Windows only)
```bash
# Build first
go build -ldflags="-s -w" -o bin/kasir-api ./cmd/api/

# Compress (Linux/Windows)
upx --best --lzma bin/kasir-api
```
**Expected size:** ~5-6M (60% reduction)

**Note:** UPX doesn't work reliably on macOS. Use Docker for cross-compilation:
```bash
# Build for Linux
docker run --rm -v "$PWD":/app -w /app golang:1.25 \
  go build -ldflags="-s -w" -o bin/kasir-api-linux ./cmd/api/

# Compress
upx --best --lzma bin/kasir-api-linux
```

### 3. Alternative Dependencies
- Replace `pgx/v5` with `lib/pq` (smaller, less features)
- Replace `koanf` with `os.Getenv()` (no .env support)
- Replace `goose` with custom migration runner

### 4. Build for Specific Platform
```bash
# Build for Linux (smaller than macOS)
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/kasir-api ./cmd/api/
```

