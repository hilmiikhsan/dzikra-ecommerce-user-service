-- +goose Up
-- +goose StatementBegin
ALTER TABLE roles ADD COLUMN static BOOLEAN NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE roles DROP COLUMN static;
-- +goose StatementEnd
