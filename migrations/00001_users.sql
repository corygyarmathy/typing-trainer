-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_name TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
);

CREATE INDEX users_created_at_idx ON users (created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
