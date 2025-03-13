-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_roles ADD COLUMN deleted_at TIMESTAMPTZ;
ALTER TABLE role_permissions ADD COLUMN deleted_at TIMESTAMPTZ;
ALTER TABLE role_app_permissions ADD COLUMN deleted_at TIMESTAMPTZ;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_roles DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE role_permissions DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE role_app_permissions DROP COLUMN IF EXISTS deleted_at;
-- +goose StatementEnd
