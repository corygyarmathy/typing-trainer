# ADR 0003: Modular monolith over microservices

- **Status:** Accepted
- **Date:** 2026-06-15

## Context

The application has clear sub-domains (`auth`, `corpus`, `progress`, `session`, `adaptive engine`) and a single consumer (the TUI client). I had to decide between a single service organised internally by domain, and multiple services each owning one domain.

## Decision

Single Go binary. Domain boundaries are enforced via package structure under `internal/`, not via process boundaries. Each domain package is independently testable and could be extracted to its own service later if scale demanded it.

## Consequences

**Positive**

- Single deployment artifact; no service-mesh, no inter-service auth, no distributed tracing required to debug a request.
- Transactions span domains where needed (e.g. registering a user creates a progress record atomically) without sagas or eventual consistency dance.
- Faster iteration for a single-author project.

**Negative**

- Domain boundaries are enforced by convention, not by network. A determined developer could create a cross-domain coupling that wouldn't be possible across services. Linting and code review mitigate this.
- A bug in one domain can take down all of them.

## Alternatives considered

- **Microservices per domain.** Rejected: unnecessary complexity at this scale. The operational complexity buys nothing because traffic patterns and failure domains don't actually differ across these contexts.
- **Two services (API + adaptive engine).** Rejected: the adaptive engine is pure-function code; running it as a separate service adds network latency and serialisation overhead with no benefit over in-process calls.
