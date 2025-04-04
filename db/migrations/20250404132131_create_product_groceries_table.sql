-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS product_groceries (
    id SERIAL PRIMARY KEY,
    min_buy INTEGER NOT NULL DEFAULT 0,
    discount BIGINT NOT NULL,
    product_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_product_groceries_product_id ON product_groceries(product_id);
CREATE INDEX idx_product_groceries_active  ON product_groceries(deleted_at) WHERE deleted_at IS NULL;

ALTER TABLE product_groceries ADD CONSTRAINT fk_product_groceries_product_id FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE ON UPDATE CASCADE;

-- Fungsi untuk memperbarui kolom `updated_at`
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW(); -- Simpan timestamp langsung dengan tipe TIMESTAMPTZ
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk tabel `product_groceries`
CREATE TRIGGER set_updated_at_product_groceries
BEFORE UPDATE ON product_groceries
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_product_groceries_active;
DROP TRIGGER IF EXISTS set_updated_at_product_groceries ON product_groceries;
DROP TABLE IF EXISTS product_groceries;
-- +goose StatementEnd
