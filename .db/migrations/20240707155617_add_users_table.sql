-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id text PRIMARY KEY,
    display_name text NOT NULL,
    email text NOT NULL,
    access_token text NOT NULL,
    refresh_token text NOT NULL,
    token_type text NOT NULL,
    expires_at timestamp NOT NULL,
    scope text NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
