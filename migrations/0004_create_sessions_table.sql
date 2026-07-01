-- +goose Up
-- +goose StatementBegin
-- One row per completed lesson. Queryable per-attempt history.
CREATE TABLE sessions (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    wpm           INT NOT NULL,
    accuracy      NUMERIC NOT NULL CHECK (accuracy >= 0 AND accuracy <= 1),
    completed_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX sessions_completed_at_idx ON sessions (completed_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sessions;
-- +goose StatementEnd
