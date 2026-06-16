# ADR 0005: sqlc for queries, goose for migrations

**Status:** Proposed
**Date:** 2026-06-15

## Context

Go has several options for database access (raw `database/sql`, `pgx` directly, query builders like `squirrel`, full ORMs like `GORM`/`ent`) and several migration tools (`golang-migrate`, `goose`, `atlas`).

## Decision

- **`sqlc`** for type-safe Go code generated from raw SQL queries.
- **`goose`** for migrations (numbered `.sql` files, applied via library at startup with an advisory lock).
- **`pgx/v5`** as the underlying driver.

## Consequences

**Positive**

- Queries are written in SQL - the language the database actually speaks. No leaky ORM abstraction, no learning a query DSL, no surprises about what SQL hits the database.
- `sqlc` generates Go types that match the SQL exactly, caught at compile time.
- `goose` migrations are plain SQL with a tiny annotation header. Familiar for anyone who's done schema migrations before.
- Advisory locking prevents race conditions when multiple instances boot against the same database (not relevant now, good hygiene later).

**Negative**

- `sqlc` can't generate everything - dynamic queries (variable `WHERE` clauses) need to be hand-written. We will use `pgx` directly for those.
- Two generated-code artifacts to keep in sync (`sqlc` and `oapi-codegen`).

## Alternatives considered

- **`GORM` / `ent`.** Rejected. `GORM` is not idiomatic and reviewers may interpret it as a sign of not having outgrown the tutorial phase. `ent` is more defensible but the schema-as-Go-code approach is heavyweight for what we need.
- **squirrel query builder.** Reasonable choice; rejected because `sqlc` gives stronger compile-time guarantees and is the current idiomatic default in modern Go backend projects.
- **`golang-migrate` over `goose`.** Both are fine. `goose` chosen for familiarity to the author.
