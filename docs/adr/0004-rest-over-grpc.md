# ADR 0004: REST + OpenAPI over gRPC

- **Status:** Proposed
- **Date:** 2026-06-15

## Context

The TUI client is the only consumer. Both REST/JSON and gRPC are credible choices in Go, with mature tooling for each.

## Decision

<!-- TODO: `oapi-codegen` is rolling out support for 3.1 - check the status  -->

REST with an OpenAPI 3.0 (not 3.1 as not yet supported by `oapi-codegen`) spec as the source of truth. Server interfaces are generated from the spec with `oapi-codegen`; handlers implement the generated interfaces.

## Consequences

**Positive**

- Universally understood. A reviewer can `curl` an endpoint with zero ceremony.
- OpenAPI gives us spec-first development (write the contract, then the code) and free documentation.
- Easier to add a browser client later if desired.

**Negative**

- Slightly higher wire overhead than gRPC. Irrelevant at this scale.
- Less rigorous type contract than protobuf. Mitigated by `oapi-codegen`'s generated types being checked by the Go compiler.

## Alternatives considered

- **gRPC.** Strong in Go and has excellent streaming. Rejected because the TUI doesn't need streaming for the v1 feature set, and gRPC's universality signal is weaker than REST's in the contexts where this project will be reviewed.
- **Hand-written JSON handlers, no spec.** Rejected because a portfolio repo without an API spec misses an opportunity to demonstrate contract-first thinking.
