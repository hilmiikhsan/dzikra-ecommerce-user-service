-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    real_price BIGINT NOT NULL,
    discount_price BIGINT NOT NULL,
    description TEXT NOT NULL,
    stock INTEGER NOT NULL DEFAULT 0,
    weight INTEGER NOT NULL DEFAULT 0,
    product_category_id INTEGER NOT NULL,
    product_sub_category_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_products_name ON products(name);
CREATE INDEX idx_products_product_category_id ON products(product_category_id);
CREATE INDEX idx_products_product_sub_category_id ON products(product_sub_category_id);
CREATE INDEX idx_products_active  ON products(deleted_at) WHERE deleted_at IS NULL;

ALTER TABLE products ADD CONSTRAINT fk_products_product_category_id FOREIGN KEY (product_category_id) REFERENCES product_categories(id) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE products ADD CONSTRAINT fk_products_product_sub_category_id FOREIGN KEY (product_sub_category_id) REFERENCES product_sub_categories(id) ON DELETE RESTRICT ON UPDATE CASCADE;

-- Fungsi untuk memperbarui kolom `updated_at`
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW(); -- Simpan timestamp langsung dengan tipe TIMESTAMPTZ
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk tabel `products`
CREATE TRIGGER set_updated_at_products
BEFORE UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_products_name;
DROP INDEX IF EXISTS idx_products_active;
DROP TRIGGER IF EXISTS set_updated_at_products ON products;
DROP TABLE IF EXISTS products;
-- +goose StatementEnd
