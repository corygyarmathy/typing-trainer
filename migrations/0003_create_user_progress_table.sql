-- +goose Up
-- +goose StatementBegin
-- One row per user, created atomically with the user at registration. The
-- competency document maps 1:1 to the engine's CompetencyState (ADR 0009);
-- the only access pattern is load-whole / write-whole per user, so it stays a
-- single JSONB document rather than normalised rows. See docs/schema.md.
CREATE TABLE user_progress (
    user_id       UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    competency    JSONB NOT NULL,-- the engine competency document
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_progress;
-- +goose StatementEnd
