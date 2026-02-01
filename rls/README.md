# Row Level Security (RLS)

This directory contains SQL scripts for managing Row Level Security policies.

## What is RLS?

Row Level Security is a **standard PostgreSQL feature** (since 9.5) that allows you to control which rows users can access in a table based on policies.

## Files

- `enable_rls.sql` - Enable RLS and create policies
- `disable_rls.sql` - Disable RLS and drop policies

## Usage

```bash
# Enable RLS
make rls-on

# Disable RLS
make rls-off
```

## Default Policies

When enabled, the following policies are created:

### Public Read Access
- Anyone can SELECT (read) from products and categories
- Useful for public-facing data

### Write Access (Choose One)

**For Supabase:**
Uncomment the `auth_write_*` policies in `enable_rls.sql`:
```sql
CREATE POLICY "auth_write_products" ON products 
    FOR ALL 
    USING (auth.role() = 'authenticated')
    WITH CHECK (auth.role() = 'authenticated');
```

**For Standard PostgreSQL:**
Uncomment the `app_write_*` policies and replace `'app_user'` with your application user:
```sql
CREATE POLICY "app_write_products" ON products 
    FOR ALL 
    USING (current_user = 'app_user')
    WITH CHECK (current_user = 'app_user');
```

## Important Notes

### Service Role Bypass
When your API connects with a **service role** (superuser or table owner), RLS policies are **bypassed**. This is the default behavior for server-side applications.

### When RLS is Enforced
RLS is enforced when:
- Using Supabase client libraries (supabase-js) with anon/authenticated keys
- Connecting with non-privileged database users
- Direct client-to-database connections

### Customizing Policies

Edit `enable_rls.sql` to customize policies for your use case:

```sql
-- Example: User can only see their own data
CREATE POLICY "user_own_data" ON products 
    FOR ALL 
    USING (user_id = current_user);

-- Example: Admin can see everything
CREATE POLICY "admin_all" ON products 
    FOR ALL 
    USING (current_user IN (SELECT username FROM admins));
```

## Learn More

- [PostgreSQL RLS Documentation](https://www.postgresql.org/docs/current/ddl-rowsecurity.html)
- [Supabase RLS Guide](https://supabase.com/docs/guides/auth/row-level-security)
