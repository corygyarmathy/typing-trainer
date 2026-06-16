# Architecture

> Read [`adr/`](adr/) first if you want to know _why_ - this document covers _what_.

## System overview

```
┌──────────────┐         HTTP / JSON         ┌────────────────────────┐
│  TUI client  │ ──────────────────────────► │  typing-trainer server │
│ (Bubble Tea) │                             │      (Go monolith)     │
└──────────────┘                             └──────────┬─────────────┘
                                                        │
                                                        │ pgx
                                                        ▼
                                              ┌──────────────────┐
                                              │   PostgreSQL 18  │
                                              └──────────────────┘
```

The server is a single binary. The TUI client is a separate binary that consumes the API.

Three executables sharing the domain code:

- `cmd/server/` - REST API, adaptive lesson engine, interacts w/ DB via SQL
- `cmd/tui/` - standalone Bubble Tea client, runs locally, talks to a configured API URL
- `cmd/sshd/` - wish-based SSH server that serves the _same_ Bubble Tea program as `cmd/tui/`, but to remote connections

## Bounded contexts

The server is structured as a modular monolith. Four bounded contexts, each owning its own types, business logic, and persistence:

| Context    | Responsibility                                         |
| ---------- | ------------------------------------------------------ |
| `auth`     | Registration, login, JWT issue and validation          |
| `corpus`   | Ngram lists, word generators, lesson source data       |
| `progress` | Per-user per-key and per-ngram competency state        |
| `session`  | Records of completed typing sessions and their results |

The `adaptive` package sits above these contexts, consuming `corpus` and `progress` to produce lessons. It contains the interesting domain logic and is deliberately I/O-free.

## Layering inside a context

Each context follows a strict three-layer shape:

- **handler.go** - HTTP transport. Parses the request, calls into the service, formats the response. Knows about HTTP; knows nothing about SQL.
- **service.go** - business logic. Orchestrates one or more repositories, manages transactions, enforces invariants. Knows about domain types; knows nothing about HTTP or SQL specifics.
- **repository.go** - persistence. Translates between SQL rows and domain types. Exposes an interface; the concrete implementation uses sqlc-generated code.

Dependencies point downward. Services do not depend on handlers; repositories do not depend on services. Wiring happens in `cmd/server/main.go`.

## Cross-cutting infrastructure

`internal/platform/` holds the load-bearing infrastructure: config loading, the database pool, logging setup, observability. These packages are used by everything but depend on nothing in the bounded contexts.

## Request flow

A typical authenticated request (`GET /api/v1/lessons/next`) flows through:

1. chi router dispatches to the registered handler
2. Recovery / RequestID / Logging middleware wrap the handler
3. Auth middleware validates the JWT and injects the user ID into the context
4. The progress handler parses the request and calls `progress.Service.NextLesson(ctx, userID)`
5. The service loads the user's competency state from `progress.Repository`
6. The service calls `adaptive.Engine.NextLesson(state, corpus)` - pure
7. The service returns the result up through the handler, which JSON-encodes it

## What's not in v1

- Microservices (modular monolith is the right size)
- Message queues (no async work yet)
- Redis cache (no measured need)
- Kubernetes (Docker on a single homelab host is fine)
- GraphQL (REST is the right fit for a TUI consumer)
- OAuth (JWT registration is enough for v1)
