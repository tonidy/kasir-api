# Database Migrations

This directory contains database migrations managed by [goose](https://github.com/pressly/goose).

## Migration Files

Migrations use the format: `{version}_{description}.sql`

Each file contains both Up and Down migrations:
```sql
-- +goose Up
-- SQL for applying the migration

-- +goose Down
-- SQL for rolling back the migration
```

## Row Level Security (RLS)

RLS is enabled on all tables with the following policies:

### Products & Categories
- **Public read access**: Anyone can SELECT
- **Authenticated write access**: Only authenticated users can INSERT/UPDATE/DELETE

### Service Role Connection
When connecting with the service role (as this API does), RLS policies are **bypassed**. This is intentional for server-side operations.

### Client-side Access
If you expose Supabase directly to clients (e.g., via supabase-js), RLS policies will be enforced based on the user's authentication status.

## Running Migrations

```bash
# Run all pending migrations
make migrate

# Or manually
go run ./cmd/api migrate
```

## Creating New Migrations

```bash
# Create a new migration file
goose -dir migrations create migration_name sql
```

This will create a new file: `{timestamp}_migration_name.sql`
