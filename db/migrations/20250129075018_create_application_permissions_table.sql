-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS application_permissions (
    id UUID PRIMARY KEY,
    application_id UUID NOT NULL,
    permission_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_application_permission UNIQUE (application_id, permission_id)
);

ALTER TABLE application_permissions ADD CONSTRAINT fk_application_permissions_application_id FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE application_permissions ADD CONSTRAINT fk_application_permissions_permission_id FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE ON UPDATE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS application_permissions;
-- +goose StatementEnd
