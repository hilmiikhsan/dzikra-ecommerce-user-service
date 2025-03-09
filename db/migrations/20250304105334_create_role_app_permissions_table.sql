-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS role_app_permissions (
    id UUID PRIMARY KEY,
    role_id UUID NOT NULL,
    app_permission_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE role_app_permissions ADD CONSTRAINT fk_role_app_permissions_role_id FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE role_app_permissions ADD CONSTRAINT fk_role_app_permissions_app_permission_id FOREIGN KEY (app_permission_id) REFERENCES application_permissions(id) ON DELETE CASCADE ON UPDATE CASCADE;

-- Fungsi untuk memperbarui kolom updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk tabel role_app_permissions
CREATE TRIGGER set_updated_at_role_app_permissions
BEFORE UPDATE ON role_app_permissions
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS role_app_permissions;
-- +goose StatementEnd
