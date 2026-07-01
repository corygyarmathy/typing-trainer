# ADR 0018: PostgreSQL as the datastore

- **Status:** Accepted
- **Date:** 2026-07-01

## Context

The application persists relational data: identities, refresh tokens, per-user progress, and typing sessions, with foreign-key relationships between them (see ADRs 0010, 0012, 0014). Several later ADRs already assume a relational store and Postgres specifically (ADR 0005 selects `sqlc`/`goose`/`pgx` _for_ Postgres), but the foundational choice of datastore was never recorded on its own. This ADR states it explicitly so the tooling decisions that depend on it have a premise to reference.

The workload is modest: a single-node modular monolith (ADR 0003) with a small number of concurrent users. The data model is squarely relational, so a document or key-value store is a poor fit; the real question is _which_ SQL database.

## Decision

**PostgreSQL.** The two primary drivers are familiarity and convention:

- I already know Postgres well, so there is no learning cost and I'm confident it covers the use case - the same velocity rationale behind the tooling in ADR 0005.
- Postgres is the idiomatic default for a Go backend. Reviewers of a portfolio project expect it, the surrounding ecosystem (`pgx`, `sqlc`, `goose`, testcontainers) is built around it first, and choosing it keeps the project on the well-trodden path.

Features like `jsonb`, native array/`uuid`/`timestamptz` types, and advisory locks (used to serialise migrations on boot, ADR 0005) are welcome and used where convenient, but they were not the reason for the choice.

## Consequences

**Positive**

- Foundational premise for ADR 0005 (`sqlc`/`goose`/`pgx`) and the identity/progress schemas is now explicit rather than assumed.
- Mature, well-documented, first-class ecosystem support in Go; excellent local/CI story via Docker and testcontainers.
- Native types and advisory locking are available when needed without reaching for extensions.

**Negative**

- Operationally heavier than an embedded database: requires a running server (managed via Docker, ADR 0006) rather than a single file. For the current single-node scale this is more infrastructure than the load strictly demands.
- A network dependency in the boot path - the app cannot start without a reachable database, whereas an embedded store would remove that failure mode.

## Alternatives considered

- **SQLite.** Not seriously considered, though defensible: embedded, zero-ops, and comfortably sufficient for a single-node typing app. Rejected because the familiarity and idiomatic-default advantages of Postgres outweighed SQLite's operational simplicity, and because running a real database server keeps the project shaped like a production service rather than a prototype.
- **MySQL / MariaDB.** The other mainstream server RDBMS. Rejected: no familiarity advantage over Postgres, weaker native type support (`jsonb`, arrays), and Postgres is the more common idiomatic pairing in modern Go backends.
- **A document store (e.g. MongoDB).** Rejected: the data is relational with cross-entity foreign keys and transactional invariants (e.g. registering a user atomically creates a progress record, ADR 0003). A relational database models this directly; a document store would fight it. Additionally, I lack familiarity with NoSQL / document-oriented databases.
