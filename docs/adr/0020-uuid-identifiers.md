# ADR 0020: UUID identifiers for entities and requests

- **Status:** Proposed
- **Date:** 2026-07-02
- **Related Artefacts:**
  - Constrains: `migrations/*.sql`, `internal/api/middleware.go`

## Context

Every persisted entity needs a primary key, and every request needs a
correlation ID for logging. Both the API process and Postgres mint identifiers,
and they do so independently - the API cannot ask the database for an ID before
it has a row, and vice versa. The choice is between coordinated identifiers
(e.g. a sequential `bigint` allocated by the database) and self-assigned random
identifiers that need no coordination.

## Decision

Use random UUIDs (RFC 4122 version 4) for both:

- **Entity primary keys** are `UUID PRIMARY KEY DEFAULT gen_random_uuid()`.
- **Request IDs** are minted in middleware via `github.com/google/uuid`.

Uniqueness of primary keys is **enforced by the `PRIMARY KEY` constraint**, not
assumed. A collision (astronomically improbable at v4's 122 random bits) surfaces
as a unique-violation on `INSERT`, never as silent data corruption. The
randomness exists to remove the need for a central ID allocator, not to
guarantee correctness on its own.

Request IDs have no uniqueness constraint and need none: they only correlate log
lines over a practical window, and a collision would merely share a trace ID
between two requests - harmless.

## Consequences

**Positive**

- No coordination: API and database can both create IDs without a round trip.
- IDs are non-sequential, so they don't leak row counts or creation order to
  clients.
- Correctness of keys rests on a database constraint, which is where it belongs.

**Negative**

- v4 UUIDs are random, so inserts scatter across the primary-key B-tree, giving
  worse index locality than a sequential key. Negligible at this project's
  scale, but real.
- 16 bytes per key vs 8 for a `bigint`.

## Alternatives considered

- **Sequential `bigint`.** Best index locality and smallest key, but leaks
  counts/ordering and requires the database to allocate every ID. Rejected: the
  coordination and information-leak costs outweigh the storage saving here.

## Future option (revisit trigger: insert-heavy tables show index bloat)

Switch entity keys to **UUIDv7** (time-ordered): keeps the no-coordination and
non-guessability properties while restoring index locality, since v7 values sort
roughly by creation time. Postgres 18 exposes `uuidv7()`; until then it can be
generated application-side. Revisit if a hot table's write throughput or index
size becomes a concern.
