# ADR 0010: Unified identity with multiple auth methods, JWT as the single API credential

- **Status:** Proposed
- **Date:** 2026-06-16
- **Related Artefacts**:
  - Supplemented by:
    - [ADR 0014](/docs/adr/0014-engine-as-library-state-follows-identity.md)
    - [ADR 0015](/docs/adr/0015-access-and-refresh-tokens.md)
  - See also: [ADR 0016](/docs/adr/0016-asymmetric-jwt-signing.md)

## Context

The system has two ways a user arrives: the HTTP API, where they register and log in with email and password, and the SSH-hosted TUI, where they are identified by their SSH public key ([ADR 0008](/docs/adr/0008-ssh-public-key-authentication.md), with an anonymous fallback when no key is offered.

Left unreconciled, that is two identity systems. I need a single identity model that both entry points map onto, and I need to decide how the SSH layer authenticates to the API once it has identified the user.

## Decision

**Decouple identity from credentials.** One `users` table holds identity. A separate `auth_credentials` table holds zero or more credentials per user:

- `kind = 'password'` - `identifier` is the (lowercase-normalised) email, `secret` is the argon2id hash.
- `kind = 'ssh_key'` - `identifier` is the public key fingerprint, `secret` is null.

A password registration creates a user plus a `password` credential. An incoming SSH connection looks up the offered key's fingerprint; if found, that is the user; if not, a user plus an `ssh_key` credential is auto-created (low-friction registration, per [ADR 0008](/docs/adr/0008-ssh-public-key-authentication.md). The same user can hold both kinds later without a schema change.

**Anonymous SSH sessions are never persisted.** When no key is offered, the session runs against an ephemeral in-memory user. No `users`, `auth_credentials`, `user_progress`, or `sessions` rows are written. This is what makes the public SSH surface safe against table-flooding: an unauthenticated connector cannot create database rows.

**JWT is the single API credential, acquired two ways.** Everything that talks to the API authenticates with a short-lived JWT access token. There are two ways to obtain one:

1. The public `POST /api/v1/auth/login` endpoint, for password users.
2. An internal, sshd-only token-exchange endpoint that mints a JWT for an SSH-resolved user. Full detail in [ADR 0016](/docs/adr/0016-asymmetric-jwt-signing.md).

The SSH layer authenticates the user via the SSH protocol (pubkey), exchanges that identity for a JWT, and from then on the TUI-over-SSH calls the API over HTTP exactly like the standalone TUI does. One client codebase, one API auth model, two transports.

**Token lifetime and renewal.** Access tokens are short-lived (15 minutes). Password-authenticated clients receive a companion refresh token (30 days, rotating) and silently renew on a `401`. The SSH path does not use refresh tokens: sshd holds the user's identity for the life of the connection and re-mints a fresh access token as needed. Full policy in [ADR 0015](/docs/adr/0015-access-and-refresh-tokens.md).

**Anonymous sessions never call the API.** Anonymous users run the engine in-process against ephemeral or local state and make no HTTP calls. Only identified users (password or SSH-key) route through the API. Full detail in [ADR 0014](/docs/adr/0014-engine-as-library-state-follows-identity.md).

## Consequences

**Positive**

- One identity, one `users` table; the API never needs to know which transport a request originated from - it sees a JWT either way.
- The schema cleanly represents password users, SSH-key users, and (by their absence) anonymous users.
- The SSH-hosted TUI and the standalone TUI share their entire client codebase; the only difference is how each acquires its initial token.
- Linking a second credential to an existing user (e.g. an SSH user later setting a password) needs no migration.

**Negative**

- More moving parts than putting `email`/`password_hash` directly on `users`. Mitigated by the join being trivial and the extension being worth it.
- The internal token-exchange endpoint is a trusted path that must be kept off the public network. A misconfiguration that exposed it would let anyone mint a token for any SSH user. Mitigated by binding it to localhost / an internal interface and treating it as a documented security boundary.
- Two token-acquisition paths to test instead of one.

## Alternatives considered

- **Nullable `email`/`password_hash`/`ssh_key_fingerprint` columns on `users`.** Simpler - no join - but every new auth method adds a column, the columns are mutually-exclusive-but-not-enforced, and it does not model "one user, several credentials." Rejected as the lower-ceiling option; the credentials table is barely more work.
- **Two separate user populations (API users and SSH users never unified).** Rejected: a reviewer who logs in via the API and then tries the SSH demo would be two different people with two progress histories, which is confusing and undermines the "same system, two front doors" story.
- **Issue JWTs to SSH users too, via the public login endpoint.** Rejected: SSH users have no password to present there. The internal exchange endpoint keeps the public surface password-only while still funnelling everything to one JWT model.
- **No JWT for the SSH path; have sshd call services in-process.** Rejected: it would fork the client into an HTTP version and an in-process version, breaking the "the SSH demo proves the same API works end to end" claim.
