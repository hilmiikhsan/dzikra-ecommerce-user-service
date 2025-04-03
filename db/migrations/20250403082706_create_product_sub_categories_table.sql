-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS product_sub_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    product_category_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_product_sub_categories_name ON product_sub_categories(name);
CREATE INDEX idx_product_sub_categories_product_category_id ON product_sub_categories(product_category_id);
CREATE INDEX idx_product_sub_categories_active  ON product_sub_categories(deleted_at) WHERE deleted_at IS NULL;

ALTER TABLE product_sub_categories ADD CONSTRAINT fk_product_sub_categories_product_category_id FOREIGN KEY (product_category_id) REFERENCES product_categories(id) ON DELETE CASCADE ON UPDATE CASCADE;

-- Fungsi untuk memperbarui kolom `updated_at`
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW(); -- Simpan timestamp langsung dengan tipe TIMESTAMPTZ
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk tabel `product_sub_categories`
CREATE TRIGGER set_updated_at_product_sub_categories
BEFORE UPDATE ON product_sub_categories
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_product_sub_categories_name;
DROP INDEX IF EXISTS idx_product_sub_categories_active;
DROP TRIGGER IF EXISTS set_updated_at_product_sub_categories ON product_sub_categories;
DROP TABLE IF EXISTS product_sub_categories;
-- +goose StatementEnd
