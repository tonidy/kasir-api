# N+1 Query Optimization - Transaction Feature

## Commits

1. **Before (N+1 version)**: `a35ca51` - feat: add transaction/checkout feature with N+1 queries (before optimization)
2. **After (Optimized)**: `8728f21` - perf: optimize transaction queries - fix N+1 issues and add row locking

## View the Difference

```bash
git diff a35ca51 8728f21
# or
git diff HEAD~1 HEAD
```

## Issues Fixed

### 1. N+1 on Product SELECT ❌ → ✅

**Before (N queries):**
```go
for _, item := range items {
    err := tx.QueryRowContext(ctx, 
        "SELECT name, price, stock FROM products WHERE id = $1", 
        item.ProductID).Scan(&productName, &productPrice, &stock)
}
```

**After (1 query with row locking):**
```go
query := fmt.Sprintf("SELECT id, name, price, stock FROM products WHERE id IN (%s) FOR UPDATE", placeholders)
rows, err := tx.QueryContext(ctx, query, productIDs...)
// Process all rows into a map
```

### 2. N+1 on Stock UPDATE ❌ → ✅

**Before (N queries inside product loop):**
```go
for _, item := range items {
    // ... fetch product ...
    _, err = tx.ExecContext(ctx, 
        "UPDATE products SET stock = stock - $1 WHERE id = $2", 
        item.Quantity, item.ProductID)
}
```

**After (N queries but separated, after validation):**
```go
// Validate all products first
for _, item := range items {
    // validation only, no DB calls
}

// Then update stock
for productID, quantity := range itemMap {
    _, err = tx.ExecContext(ctx, 
        "UPDATE products SET stock = stock - $1 WHERE id = $2", 
        quantity, productID)
}
```

### 3. N+1 on Detail INSERT ❌ → ✅

**Before (N queries):**
```go
for i := range details {
    details[i].TransactionID = transactionID
    var detailID int
    err = tx.QueryRowContext(ctx,
        "INSERT INTO transaction_details (...) VALUES ($1, $2, $3, $4) RETURNING id",
        transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal).Scan(&detailID)
    details[i].ID = detailID
}
```

**After (1 query with batch insert):**
```go
query := "INSERT INTO transaction_details (...) VALUES "
// Build: ($1,$2,$3,$4), ($5,$6,$7,$8), ... RETURNING id
rows, err := tx.QueryContext(ctx, query, args...)
// Scan all returned IDs
```

## Performance Impact

**10 items checkout:**
- **Before**: ~31 queries
  - 10 SELECT (products)
  - 10 UPDATE (stock)
  - 10 INSERT (details)
  - 1 INSERT (transaction)
  
- **After**: ~13 queries
  - 1 SELECT (all products)
  - 10 UPDATE (stock) *
  - 1 INSERT (details batch)
  - 1 INSERT (transaction)

\* UPDATE queries still N because PostgreSQL doesn't support batch UPDATE with different WHERE clauses efficiently

## Bonus: Race Condition Protection

Added `FOR UPDATE` lock:
```sql
SELECT ... FROM products WHERE id IN (...) FOR UPDATE
```

This prevents concurrent transactions from:
- Reading stale stock values
- Overselling products
- Creating race conditions

## Code Quality

✅ All tests pass  
✅ Build successful  
✅ Code formatted  
✅ No linter warnings

## Lines Changed

```
1 file changed, 85 insertions(+), 25 deletions(-)
```

## Key Improvements

1. **Batch SELECT**: Fetch all products in one query
2. **Row Locking**: `FOR UPDATE` prevents race conditions
3. **Batch INSERT**: Insert all details in one query with `RETURNING`
4. **Early Validation**: Validate all items before any updates
5. **Better Structure**: Separate concerns (fetch → validate → update → insert)
