# Product Filtering Feature

## Summary
Added query parameter filtering to the products API endpoint to support filtering by name and active status.

## Changes Made

### 1. Model Updates
- **File**: `internal/model/product.go`
- Added `Active bool` field to Product struct

### 2. Repository Layer
- **File**: `internal/repository/interfaces.go`
  - Added `FindByFilters(ctx context.Context, name string, active *bool) ([]model.Product, error)` to ProductReader interface

- **File**: `internal/repository/memory/product.go`
  - Implemented `FindByFilters` method with case-insensitive name matching using `strings.Contains`
  - Filters by name (partial match) and/or active status

- **File**: `internal/repository/postgres/product.go`
  - Implemented `FindByFilters` method with SQL ILIKE for case-insensitive name matching
  - Updated all queries to include the `active` column
  - Updated `Create` and `Update` methods to handle the active field

### 3. Service Layer
- **File**: `internal/service/product.go`
- Added `GetByFilters(ctx context.Context, name string, active *bool) ([]model.Product, error)` method

### 4. Handler Layer
- **File**: `internal/handler/product.go`
  - Updated `GetAll` handler to parse query parameters (`name` and `active`)
  - Routes to `GetByFilters` when filters are present, otherwise uses `GetAll`
  - Added `GetByFilters` to ProductService interface

- **File**: `internal/dto/product.go`
  - Added `Active bool` field to ProductResponse

### 5. Tests
- **File**: `internal/handler/product_test.go`
- Updated mockProductService to implement GetByFilters method

### 6. Database Migration
- **File**: `database/migrations/00003_add_active_to_products.sql`
- Added migration to add `active` column to products table with default value `true`

## API Usage

### Get All Products
```bash
GET /api/products
```

### Filter by Name (case-insensitive, partial match)
```bash
GET /api/products?name=indom
```

### Filter by Active Status
```bash
GET /api/products?active=true
GET /api/products?active=false
```

### Combined Filters
```bash
GET /api/products?name=indom&active=true
```

## Example Response
```json
[
  {
    "id": 1,
    "name": "Indomie Goreng",
    "price": 3500,
    "stock": 100,
    "active": true,
    "category": {
      "id": 1,
      "name": "Food",
      "description": "Food and snacks"
    }
  }
]
```

## Testing
All existing tests pass. The feature has been manually tested with:
- Name filtering (partial, case-insensitive)
- Active status filtering (true/false)
- Combined name and active filtering
- No filters (returns all products)

## Migration
Run the migration to add the active column:
```bash
./bin/kasir-api migrate
```

Or using make:
```bash
make migrate
```
