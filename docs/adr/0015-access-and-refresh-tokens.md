# ADR 0015: Short-lived access tokens with rotating refresh tokens

- **Status:** Proposed
- **Date:** 2026-06-22
- **Related Artefacts:**
  - Supplements: [ADR 0010](/docs/adr/0010-unified-identity-and-jwt.md)
  - References: [ADR 0014](/docs/adr/0014-engine-as-library-state-follows-identity.md)

## Context

[ADR 0010](/docs/adr/0010-unified-identity-and-jwt.md) established the JWT as the single API credential but left its lifetime unspecified; the draft OpenAPI spec set `expires_in` to 900 seconds (15 minutes) with no way to renew. Fifteen minutes is shorter than a focused practice session plus reading the results, and a standalone TUI talking to a remote API would hit a `401` mid-use with no recovery path. The system needs a token-lifetime policy and a renewal mechanism.

## Decision

Adopt the standard access / refresh split.

- **Access token.** Short-lived JWT (15 minutes), sent as `Authorization: Bearer`. Stateless; verified on every request.
- **Refresh token.** Long-lived (30 days), opaque, **rotating**. Exchanged for a new access token _and a new refresh token_ at a new `POST /api/v1/auth/refresh`. Rotation means each refresh invalidates the previous refresh token, so a leaked-and-replayed refresh token is detectable and bounded.
- `POST /api/v1/auth/register` and `POST /api/v1/auth/login` return both tokens.
- **Standalone TUI (password).** Stores the refresh token in its local state file (`$XDG_STATE_HOME/typist/`, per [ADR 0014](/docs/adr/0014-engine-as-library-state-follows-identity.md)) and silently refreshes on a `401`.
- **SSH path.** Needs no refresh token. sshd holds the user's identity for the life of the connection ([ADR 0010](/docs/adr/0010-unified-identity-and-jwt.md)) and mints a fresh access token per connection; if one expires mid-session it re-mints. The access-token lifetime is set comfortably longer than a single lesson, so this is rare.

## Consequences

**Positive**

- Sessions survive well beyond 15 minutes without lengthening the access token, whose short life limits the blast radius of a leak.
- Rotation gives basic refresh-token-theft detection for the cost of one table.
- The split is the idiomatic pattern reviewers expect; demonstrating it is a portfolio plus.

**Negative**

- A refresh token is server-side state: a `refresh_tokens` table (hashed token, user, expiry, rotation lineage) and a revocation check. Acceptable - one small table and one endpoint.
- The client must implement transparent refresh-on-`401`. Contained in the API-client layer and shared by both transports.

## Alternatives considered

- **A single moderate-lived access token (12-24h), no refresh.** Simpler, but no revocation and a longer exposure window for a leaked token; a weaker signal. Rejected, though it remains the honest fallback if the refresh table proves not worth its weight.
- **15-minute access token, no refresh (the draft).** Rejected: breaks an ordinary practice session.
- **OAuth2 / a third-party identity provider.** Rejected per [ADR 0010](/docs/adr/0010-unified-identity-and-jwt.md)'s scope - no third-party login in v1.
