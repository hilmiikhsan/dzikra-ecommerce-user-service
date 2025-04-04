-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS product_variants (
    id SERIAL PRIMARY KEY,
    variant_sub_name VARCHAR(255) NOT NULL,
    variant_stock INTEGER NOT NULL DEFAULT 0,
    variant_weight DECIMAL(10, 2) NOT NULL DEFAULT 0,
    capital_price BIGINT NOT NULL,
    real_price BIGINT NOT NULL,
    discount_price BIGINT NOT NULL,
    product_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_product_variants_product_id ON product_variants(product_id);
CREATE INDEX idx_product_variants_active  ON product_variants(deleted_at) WHERE deleted_at IS NULL;

ALTER TABLE product_variants ADD CONSTRAINT fk_product_variants_product_id FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE ON UPDATE CASCADE;

-- Fungsi untuk memperbarui kolom `updated_at`
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW(); -- Simpan timestamp langsung dengan tipe TIMESTAMPTZ
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk tabel `product_variants`
CREATE TRIGGER set_updated_at_product_variants
BEFORE UPDATE ON product_variants
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_product_variants_active;
DROP TRIGGER IF EXISTS set_updated_at_product_variants ON product_variants;
DROP TABLE IF EXISTS product_variants;
-- +goose StatementEnd
