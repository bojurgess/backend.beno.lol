-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id text PRIMARY KEY,
    display_name text,
    email text,
    access_token text,
    token_type text,
	expires_at date,
	scope text,
	refresh_token text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users
-- +goose StatementEnd