-- +goose Up
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN username;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN username VARCHAR(100) NOT NULL UNIQUE;
-- +goose StatementEnd
