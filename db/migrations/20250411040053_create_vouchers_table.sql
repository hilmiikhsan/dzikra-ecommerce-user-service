-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS vouchers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    voucher_quota INTEGER NOT NULL DEFAULT 0,
    code VARCHAR(255) NOT NULL UNIQUE,
    discount BIGINT NOT NULL DEFAULT 0,
    start_at TIMESTAMP,
    end_at TIMESTAMP,
    voucher_type_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_vouchers_name ON vouchers(name);
CREATE INDEX idx_vouchers_active  ON vouchers(deleted_at) WHERE deleted_at IS NULL;

ALTER TABLE vouchers ADD CONSTRAINT fk_vouchers_voucher_type_id FOREIGN KEY (voucher_type_id) REFERENCES voucher_types(id) ON DELETE RESTRICT ON UPDATE CASCADE;

-- Fungsi untuk memperbarui kolom `updated_at`
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW(); -- Simpan timestamp langsung dengan tipe TIMESTAMPTZ
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk tabel `vouchers`
CREATE TRIGGER set_updated_at_vouchers
BEFORE UPDATE ON vouchers
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_vouchers_name;
DROP INDEX IF EXISTS idx_vouchers_active;
DROP TRIGGER IF EXISTS set_updated_at_vouchers ON vouchers;
DROP TABLE IF EXISTS vouchers;
-- +goose StatementEnd
