-- +goose Up
-- +goose StatementBegin
-- Buat tipe enum untuk device_type
CREATE TYPE device_type_enum AS ENUM ('ANDROID', 'IOS');

-- Buat tabel user_fcm_tokens dengan device_type menggunakan enum
CREATE TABLE IF NOT EXISTS user_fcm_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    device_id VARCHAR(255) NOT NULL,
    device_type device_type_enum NOT NULL,
    fcm_token VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Fungsi untuk memperbarui kolom updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk tabel user_fcm_tokens
CREATE TRIGGER set_updated_at_user_fcm_tokens
BEFORE UPDATE ON user_fcm_tokens
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Hapus trigger, tabel, dan tipe enum jika diperlukan
DROP TRIGGER IF EXISTS set_updated_at_user_fcm_tokens ON user_fcm_tokens;
DROP TABLE IF EXISTS user_fcm_tokens;
DROP TYPE IF EXISTS device_type_enum;
-- +goose StatementEnd
