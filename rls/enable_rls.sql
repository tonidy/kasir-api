-- Enable Row Level Security on all tables
ALTER TABLE products ENABLE ROW LEVEL SECURITY;
ALTER TABLE categories ENABLE ROW LEVEL SECURITY;

-- Public read access (anyone can SELECT)
CREATE POLICY "public_read_products" ON products 
    FOR SELECT 
    USING (true);

CREATE POLICY "public_read_categories" ON categories 
    FOR SELECT 
    USING (true);

-- For Supabase: Authenticated users can write
-- Uncomment if using Supabase with auth.role()
-- CREATE POLICY "auth_write_products" ON products 
--     FOR ALL 
--     USING (auth.role() = 'authenticated')
--     WITH CHECK (auth.role() = 'authenticated');

-- CREATE POLICY "auth_write_categories" ON categories 
--     FOR ALL 
--     USING (auth.role() = 'authenticated')
--     WITH CHECK (auth.role() = 'authenticated');

-- For standard PostgreSQL: Specific user/role can write
-- Uncomment and replace 'app_user' with your application user
-- CREATE POLICY "app_write_products" ON products 
--     FOR ALL 
--     USING (current_user = 'app_user')
--     WITH CHECK (current_user = 'app_user');

-- CREATE POLICY "app_write_categories" ON categories 
--     FOR ALL 
--     USING (current_user = 'app_user')
--     WITH CHECK (current_user = 'app_user');
