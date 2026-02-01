-- +migrate Up
ALTER TABLE products ADD COLUMN category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL;
CREATE INDEX idx_products_category_id ON products(category_id);

-- +migrate Down
DROP INDEX IF EXISTS idx_products_category_id;
ALTER TABLE products DROP COLUMN category_id;
