# Database Seeds

This directory contains seed data for development and testing.

## Seed Files

- `seed.sql` - Sample data for products and categories

## Running Seeds

```bash
# Seed the database
make seed

# Or manually
go run ./cmd/api seed
```

## Seed Data

The seed file includes:
- 4 categories (Food, Beverage, Electronics, Stationery)
- 9 sample products across different categories

Seeds use `ON CONFLICT DO NOTHING` to be idempotent - you can run them multiple times safely.

## Creating Custom Seeds

Add your SQL statements to `seed.sql`:

```sql
INSERT INTO table_name (columns) VALUES (values)
ON CONFLICT DO NOTHING;
```

The `ON CONFLICT DO NOTHING` clause ensures seeds can be run multiple times without errors.
