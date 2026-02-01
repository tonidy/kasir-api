-- Drop all RLS policies
DROP POLICY IF EXISTS "public_read_products" ON products;
DROP POLICY IF EXISTS "public_read_categories" ON categories;
DROP POLICY IF EXISTS "auth_write_products" ON products;
DROP POLICY IF EXISTS "auth_write_categories" ON categories;
DROP POLICY IF EXISTS "app_write_products" ON products;
DROP POLICY IF EXISTS "app_write_categories" ON categories;

-- Disable Row Level Security
ALTER TABLE products DISABLE ROW LEVEL SECURITY;
ALTER TABLE categories DISABLE ROW LEVEL SECURITY;
