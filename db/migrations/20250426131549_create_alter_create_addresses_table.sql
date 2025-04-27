-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS address;

CREATE TABLE IF NOT EXISTS addresses (
    id SERIAL PRIMARY KEY,
    province VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    district VARCHAR(100),
    subdistrict VARCHAR(100) NOT NULL,
    postal_code VARCHAR(100) NOT NULL,
    address VARCHAR(255) NOT NULL,
    received_name VARCHAR(100) NOT NULL,
    user_id UUID NOT NULL,
    city_vendor_id VARCHAR(100) NOT NULL,
    province_vendor_id VARCHAR(100) NOT NULL,
    subdistrict_vendor_id VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

ALTER TABLE addresses ADD CONSTRAINT fk_addresses_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE;

CREATE INDEX idx_addresses_user_id ON addresses (user_id);
CREATE INDEX idx_addresses_province_city_district ON addresses (province, city, district);
CREATE INDEX idx_addresses_postal_code ON addresses (postal_code);
CREATE INDEX idx_addresses_active ON addresses (deleted_at) WHERE deleted_at IS NULL;


-- Fungsi untuk memperbarui kolom `updated_at`
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW(); -- Simpan timestamp langsung dengan tipe TIMESTAMPTZ
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk tabel `addresses`
CREATE TRIGGER set_updated_at_addresses
BEFORE UPDATE ON addresses
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_addresses_user_id;
DROP INDEX IF EXISTS idx_addresses_province_city_district;
DROP INDEX IF EXISTS idx_addresses_postal_code;
DROP INDEX IF EXISTS idx_addresses_active;
DROP TRIGGER IF EXISTS set_updated_at_addresses ON addresses;
DROP TABLE IF EXISTS addresses;
-- +goose StatementEnd
