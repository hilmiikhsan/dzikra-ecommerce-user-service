-- +goose Up
-- +goose StatementBegin
CREATE TYPE voucher_type_enum AS ENUM ('OFFLINE_SHOP', 'ONLINE_SHOP');

CREATE TABLE IF NOT EXISTS voucher_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    type voucher_type_enum UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_voucher_types_name ON voucher_types(name);
CREATE INDEX idx_voucher_types_active ON voucher_types(deleted_at) WHERE deleted_at IS NULL;

-- Fungsi untuk memperbarui kolom updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk memperbarui kolom updated_at pada tabel voucher_types
CREATE TRIGGER set_updated_at_voucher_types
BEFORE UPDATE ON voucher_types
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_updated_at_voucher_types ON voucher_types;
DROP INDEX IF EXISTS idx_voucher_types_active;
DROP INDEX IF EXISTS idx_voucher_types_name;
DROP TABLE IF EXISTS voucher_types;
DROP TYPE IF EXISTS voucher_type_enum;
-- +goose StatementEnd
