-- +goose Up
-- +goose StatementBegin
-- Identity is decoupled from credentials (ADR 0010): one users row may hold
-- several credentials of different kinds. v1 only issues 'password' rows; the
-- 'ssh_key' kind lands with the SSH surface without a schema change.
CREATE TABLE auth_credentials (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    cred_kind     TEXT NOT NULL CHECK (cred_kind IN ('password', 'ssh_key')),
    -- For 'password': the lowercase-normalised email. For 'ssh_key': the key
    -- fingerprint. Application normalises before insert and lookup.
    identifier    TEXT NOT NULL,
    -- argon2id hash for 'password'; null for 'ssh_key'.
    secret        TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),

    -- Both the integrity rule and the login / key-resolution lookup index.
    UNIQUE (cred_kind, identifier),

    -- Enforce expected secret behaviour
    CHECK (
      (cred_kind = 'ssh_key' AND secret IS NULL)
      OR
      (cred_kind = 'password' AND secret IS NOT NULL)
    )
);

CREATE INDEX auth_credentials_created_at_idx ON auth_credentials (created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE auth_credentials;
-- +goose StatementEnd
