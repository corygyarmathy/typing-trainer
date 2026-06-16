# ADR 0001: Record architecture decisions

**Status:** Accepted
**Date:** 2026-06-15

## Context

This project makes a number of architectural and structural choices - a modular
monolith over microservices, env-var configuration over a config file, `sqlc` over
an ORM, and so on - whose _reasoning_ is not visible from the code itself. Code
shows what was decided; it rarely shows what alternatives were weighed or why
they were rejected. Without that record, the rationale is lost, settled
questions get relitigated, and a newcomer (including a future version of the
author) cannot tell a deliberate decision from an accident.

This is also a portfolio project. The decisions and their justifications are
themselves part of what the project is meant to demonstrate.

## Decision

We will record architecturally significant decisions as [Architecture Decision Records (ADRs)](https://github.com/architecture-decision-record/architecture-decision-record), following the lightweight format popularised by [Michael Nygard](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions).

Conventions:

- ADRs live in `docs/adr/` as Markdown files, one decision per file.
- Files are numbered sequentially and zero-padded: `NNNN-short-title.md`. The number is permanent and never reused.
- Each ADR has the sections **Status**, **Context**, **Decision**, and **Consequences**, optionally with an **Alternatives considered** section where the trade-offs are worth recording explicitly.
- Status follows a simple lifecycle: `Proposed` → `Accepted`, and later `Deprecated` or `Superseded by NNNN` if circumstances change.
- ADRs are immutable once accepted. A decision is not edited to reflect a change of mind; instead a new ADR is written that supersedes the old one, and the old one's status is updated to point at it. The history of the thinking is part of the value.

A decision is "architecturally significant" - and therefore worth an ADR - if it is costly to reverse, affects the structure or dependencies of the system, constrains future choices, or is non-obvious enough that someone would reasonably ask "why was it done this way?".

## Consequences

- The reasoning behind non-obvious choices is preserved next to the code, in version control, and evolves with it.
- New contributors can read the ADRs in order to understand how the system came to look the way it does, rather than reverse-engineering intent from the code.
- There is a small, ongoing cost: significant decisions must be written down at the time they are made, while the context is still fresh.
- The numbered, immutable record makes superseded decisions auditable - it is always possible to see not just the current design but the path taken to it.
