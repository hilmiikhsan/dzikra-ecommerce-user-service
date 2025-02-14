-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS address (
    id SERIAL PRIMARY KEY,
    province VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    district VARCHAR(100) NOT NULL,
    subdistrict VARCHAR(100) NOT NULL,
    postal_code VARCHAR(100) NOT NULL,
    address TEXT NOT NULL,
    received_name VARCHAR(100) NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

ALTER TABLE address ADD CONSTRAINT fk_address_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE;

CREATE INDEX idx_address_user_id ON address (user_id);
CREATE INDEX idx_address_province_city_district ON address (province, city, district);

-- Fungsi untuk memperbarui kolom `updated_at`
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW(); -- Simpan timestamp langsung dengan tipe TIMESTAMPTZ
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk tabel `address`
CREATE TRIGGER set_updated_at_address
BEFORE UPDATE ON address
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_address_user_id;
DROP INDEX IF EXISTS idx_address_province_city_district;
DROP TRIGGER IF EXISTS set_updated_at_address ON address;
DROP TABLE IF EXISTS address;
-- +goose StatementEnd
