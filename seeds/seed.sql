-- Seed categories
INSERT INTO categories (name, description) VALUES
    ('Food', 'Food and snacks'),
    ('Beverage', 'Drinks and beverages'),
    ('Electronics', 'Electronic devices and accessories'),
    ('Stationery', 'Office and school supplies')
ON CONFLICT DO NOTHING;

-- Seed products
INSERT INTO products (name, price, stock, category_id) VALUES
    ('Indomie Goreng', 3500, 100, (SELECT id FROM categories WHERE name = 'Food' LIMIT 1)),
    ('Indomie Soto', 3500, 100, (SELECT id FROM categories WHERE name = 'Food' LIMIT 1)),
    ('Chitato', 12000, 50, (SELECT id FROM categories WHERE name = 'Food' LIMIT 1)),
    ('Coca Cola 330ml', 5000, 80, (SELECT id FROM categories WHERE name = 'Beverage' LIMIT 1)),
    ('Aqua 600ml', 3000, 120, (SELECT id FROM categories WHERE name = 'Beverage' LIMIT 1)),
    ('USB Cable Type-C', 25000, 30, (SELECT id FROM categories WHERE name = 'Electronics' LIMIT 1)),
    ('Power Bank 10000mAh', 150000, 15, (SELECT id FROM categories WHERE name = 'Electronics' LIMIT 1)),
    ('Ballpoint Pen', 2500, 200, (SELECT id FROM categories WHERE name = 'Stationery' LIMIT 1)),
    ('Notebook A5', 15000, 50, (SELECT id FROM categories WHERE name = 'Stationery' LIMIT 1))
ON CONFLICT DO NOTHING;
