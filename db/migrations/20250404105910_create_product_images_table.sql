-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS product_images (
    id SERIAL PRIMARY KEY,
    image_url VARCHAR(255) NOT NULL,
    product_id INTEGER NOT NULL,
    sort INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_product_images_product_id ON product_images(product_id);
CREATE INDEX idx_product_images_active  ON product_images(deleted_at) WHERE deleted_at IS NULL;

ALTER TABLE product_images ADD CONSTRAINT fk_product_images_product_id FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE ON UPDATE CASCADE;

-- Fungsi untuk memperbarui kolom `updated_at`
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW(); -- Simpan timestamp langsung dengan tipe TIMESTAMPTZ
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk tabel `product_images`
CREATE TRIGGER set_updated_at_product_images
BEFORE UPDATE ON product_images
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_product_images_active;
DROP TRIGGER IF EXISTS set_updated_at_product_images ON product_images;
DROP TABLE IF EXISTS product_images;
-- +goose StatementEnd
