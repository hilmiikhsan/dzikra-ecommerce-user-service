-- +goose Up
-- +goose StatementBegin

-- Buat tabel voucher_usage (sesuaikan nama tabel dengan kebutuhan Anda)
CREATE TABLE IF NOT EXISTS voucher_usage (
    id SERIAL PRIMARY KEY,
    is_use BOOLEAN NOT NULL,
    voucher_id INT NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_voucher_usage_voucher_id ON voucher_usage(voucher_id);
CREATE INDEX idx_voucher_usage_user_id ON voucher_usage(user_id);
CREATE INDEX idx_voucher_usage_active ON voucher_usage(deleted_at) WHERE deleted_at IS NULL;

ALTER TABLE voucher_usage ADD CONSTRAINT fk_voucher_usage_voucher_id FOREIGN KEY (voucher_id) REFERENCES vouchers(id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE voucher_usage ADD CONSTRAINT fk_voucher_usage_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE;

-- Fungsi untuk memperbarui kolom updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk tabel voucher_usage
CREATE TRIGGER set_updated_at_voucher_usage
BEFORE UPDATE ON voucher_usage
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_voucher_usage_voucher_id;
DROP INDEX IF EXISTS idx_voucher_usage_user_id;
DROP INDEX IF EXISTS idx_voucher_usage_active;
DROP TRIGGER IF EXISTS set_updated_at_voucher_usage ON voucher_usage;
DROP TABLE IF EXISTS voucher_usage;
DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE;
-- +goose StatementEnd
