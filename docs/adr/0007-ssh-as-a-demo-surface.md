# ADR 0007: SSH as a Demo Surface

**Status:** Proposed
**Date:** 2026-06-16

## Context

The purpose of this project is to be a portfolio project, and thus a live demo is a useful feature: reducing friction required for reviewers to test the product. While the server and client will still be buildable and deployable, I will also make the TUI client available over SSH.

## Decision

SSH'ing into the TUI client means that the reviewer can use the tool within seconds - no installs or configurations required. The TUI client will rely on [Bubble Tea](https://github.com/charmbracelet/bubbletea), and the SSH-accessible client will use the [Wish](https://github.com/charmbracelet/wish) library, which is purpose-built for serving Bubble Tea apps over SSH.

The specifics of authentication are addressed in [`ADR 0008: SSH Public Key Authentication`](/docs/adr/0008-ssh-public-key-authentication.md).

## Consequences

**Positive**

- Enables TUI-client app to have a 'live demo', low-friction accessibility.
- [Wish](https://github.com/charmbracelet/wish) handles the heavy lifting - SSH protocol, session multiplexing, terminal sizing, PTY allocation. Opportunity for demonstrating infrastructure experience: deploying a real app on a real server, accessible over the internet w/ good security practices in place.

**Negative**

- Three deployment surfaces instead of two: the API server, the SSH-TUI server, and the standalone TUI binary.
- Increases required development time to finish the project.
- Increases project complexity, dependencies.
- Increases attack surface and security risks.

## Alternatives considered

- **Local-client only.** Rejected: would add friction for reviewer and mean that I would not have a 'live demo' for the tool.
- **SSH only.** Rejected: I need to build the TUI-client anyway, and making it SSH-only is unnecessarily limiting the functionality.
