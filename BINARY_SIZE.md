# Binary Size History

This document tracks the evolution of the application's binary size and dependencies over time.

## Timeline

### Phase 1: Initial Implementation (Commit: initial)
**Date:** Early development  
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

**Architecture:**
```
internal/
├── config/      # Configuration management
├── database/    # DB connection & migrations
├── model/       # Domain entities
├── repository/  # Data access (memory + postgres)
├── service/     # Business logic
└── handler/     # HTTP handlers
```

---

## Size Comparison Table

| Phase | Date | Binary Size | Change | Dependencies | Key Features |
|-------|------|-------------|--------|--------------|--------------|
| 1 | Initial | 7.9M | - | 0 | In-memory, stdlib only |
| 2 | 2026-02-01 | 15M | +7.1M (+90%) | 23 | PostgreSQL, migrations, config |

---

## Optimization Options

If binary size becomes a concern, consider:

### 1. Build with Compression
```bash
# Strip debug symbols
go build -ldflags="-s -w" -o bin/kasir-api ./cmd/api/

# Compress with UPX
upx --best --lzma bin/kasir-api
```
**Expected size:** ~5-6M (60% reduction)

### 2. Alternative Dependencies
- Replace `pgx/v5` with `lib/pq` (smaller, less features)
- Replace `koanf` with `os.Getenv()` (no .env support)
- Replace `goose` with custom migration runner

### 3. Build Tags
```bash
# Build without PostgreSQL support
go build -tags=nomemory -o bin/kasir-api ./cmd/api/
```

---

## Analysis

### Why Size Increases?

1. **Database Drivers** - Full-featured drivers include:
   - Protocol implementations
   - Type conversions
   - Connection pooling
   - Error handling

2. **Configuration Libraries** - Include:
   - Multiple parsers (env, file, dotenv)
   - Type conversion
   - Validation

3. **Migration Tools** - Include:
   - SQL parsing
   - Version tracking
   - Rollback support

### Is It Worth It?

**Yes, for production applications:**
- Maintainability > Binary size
- Features > Minimal footprint
- Developer experience > Deployment size

**No, if:**
- Embedded systems (size-constrained)
- Edge computing (bandwidth-limited)
- Simple prototypes
- Learning projects

---

## Future Considerations

As the application grows, monitor:
- Binary size per major feature
- Dependency count
- Build time
- Deployment size

**Target thresholds:**
- ⚠️ Binary > 50M: Review dependencies
- ⚠️ Dependencies > 50: Consider consolidation
- ⚠️ Build time > 1min: Optimize build process

---

## How to Update This Document

When adding new dependencies:

1. **Build and measure:**
   ```bash
   make build
   du -sh bin/kasir-api
   ```

2. **Count dependencies:**
   ```bash
   wc -l go.mod go.sum
   ```

3. **Add new phase:**
   - Date and commit hash
   - Binary size and change
   - New dependencies with purpose
   - Features added

4. **Update comparison table**

---

## Recommendations

**Current Phase (Phase 2):**
- ✅ Binary size is acceptable for production (15M)
- ✅ Dependencies are well-justified
- ✅ Architecture supports long-term maintenance
- ✅ No optimization needed at this stage

**Monitor for Phase 3:**
- If adding observability (Prometheus, OpenTelemetry)
- If adding API documentation (Swagger/OpenAPI)
- If adding authentication (JWT, OAuth)
- If adding caching (Redis client)
