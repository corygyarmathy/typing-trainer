# ADR 0009: JSON serialised data

- **Status:** Proposed
- **Date:** 2026-06-16
- **Related Artefacts**:
  - [`ADR 0002: Golang Selection`](/docs/adr/0002-golang-selection.md).
  - [`ADR 0004: REST + OpenAPI over gRPC`](/docs/adr/0004-rest-over-grpc.md).

## Context

Adopting the client/server, and RESTful API architecture, some serialisation of data will be necessary. I will need to decide what format to select for my project.

## Decision

I will use JSON for the serialisation of data in my project.

## Consequences

**Positive**

- JSON is the most commonly used serialisation type for APIs, and therefore will have the most compatibility with existing tools and workflows.
- The format is human readable.

**Neutral**

- [OpenAPI assumes JSON or YAML](https://spec.openapis.org/oas/v3.1.0.html#format).
- `oapi-codegen` doesn't generate Gob handlers.
- While the performance of JSON isn't amazing, for this context the performance impact in comparison to other formats is irrelevant.

## Alternatives considered

- **Gob.** Rejected:
  - Gob is Go-only. The moment a Rust client, a `curl` probe, a Postman test, a JavaScript browser client, or a colleague's Python script wants to talk to the API, they can't. Even though the TUI is Go today, I don't want to weld that constraint into the contract.
  - Gob isn't introspectable. A reviewer can't `curl /api/v1/lessons/next | jq` to see what the API returns. They'd have to write a Go client first. Adds unnecessary friction.
