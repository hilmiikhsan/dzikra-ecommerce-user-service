-- +goose Up
-- +goose StatementBegin
ALTER TABLE role_permissions DROP CONSTRAINT unique_role_permissions;
CREATE UNIQUE INDEX unique_role_permissions_active
ON role_permissions(role_id, permission_id)
WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS unique_role_permissions_active;
ALTER TABLE role_permissions ADD CONSTRAINT unique_role_permissions UNIQUE (role_id, permission_id);
-- +goose StatementEnd
