# ADR 0008: SSH Public Key Authentication

**Status:** Proposed
**Date:** 2026-06-16

## Context

As part of the live demo functionality, I would like a low-friction, secure, low-complexity authentication measure.

## Decision

The SSH client will authenticate with a user's SSH public key, which uniquely identifies each host. Where the public key is not exposed, it will fallback to an ephemeral anonymous session, with a note displayed to the user as to why.

Due to significant security risks, the following controls are planned:

- Ephemeral (never-persisted) anonymous sessions
- Per-IP connection rate limiting
- A cap on concurrent sessions
- Running the sshd container as an unprivileged user with locked-down egress
- Host-key management
- Confirm a connection can't escape the Bubble Tea program into a shell (Wish handles PTY but verify)

## Consequences

**Positive**

- Authentication is (mostly) seamless to the user, assuming they expose their public key when SSH'ing.
- Each user is automatically identified and tracked without needing to implement a complex registration system.
- Public keys are a good practice method of identifying hosts, which is already built into standard developer workflows and tools.

**Negative**

- Increased system complexity compared to no authentication (i.e. just having ephemeral anonymous sessions).
- The same user may be assumed to be different users if connecting from different hosts, depending on their SSH configuration.
- A publicly-SSH-able service is an attack surface in a way that an HTTP API behind Cloudflare isn't. Requires additional security measures.

## Alternatives considered

- **Local-client only.** Rejected: would add friction for reviewer and mean that I would not have a 'live demo' for the tool.
- **SSH only.** Rejected: It would be an unnecessary restriction, as this is not a proprietary project. I need to build the client regardless, so it would not save any time.
