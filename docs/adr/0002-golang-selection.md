# ADR 0002: Golang Selection

- **Status:** Accepted
- **Date:** 2026-06-15

## Context

This is a portfolio project, completed as part of the end of the [Boot.Dev](https://boot.dev) Backend Software Development course. In the course, you learn two primary languages: Python and Go.

The purpose of this project is to ultimately be presented as part of a job application, and so therefore the context of the Perth (and Australian) job market must be considered.

## Decision

I have decided to write the project in Go. This is for the following reasons:

- Existing familiarity and knowledge of the language.
- Existing familiarity and knowledge of Go tools / ecosystem.
- Go is performant, especially when compared to interpreted languages.
- Go's ecosystem has mature tools for the problem space.
- Go's static strong type system assists with enforcing the project structure.

## Consequences

- The major effect is that all the code will be written in Go.
- Impacts all chosen tools, as I will be choosing ones that work well with Go.

## Alternatives considered

- **Java/C#.** Rejected: Different languages, but rejected for the same reasons. The major reasons for considering these languages is that they are prevalent in the Perth job market, where Go isn't as much. While I have familiarity with these languages, it is not currently at the same level as Go or Python. To learn a language and toolset on top of completing this project would delay it too much.
- **Python.** Rejected: The major reason for considering this is both my familiarity with the tool and its prevalence in the Australian job market. The main reasons for rejecting it are: I'm more familiar with and confident in the Go tools required for this project, and I've found that - due to its flexibility - Python requires the developer to enforce structure rather than the language doing so, which creates issues as the project complexity increases.
