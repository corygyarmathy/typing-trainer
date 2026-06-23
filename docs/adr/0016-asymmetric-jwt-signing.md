# ADR 0016: Asymmetric JWT signing for the internal token exchange

- **Status:** Proposed
- **Date:** 2026-06-22
- **Related Artefacts:**
  - Supplements: [ADR 0010](/docs/adr/0010-unified-identity-and-jwt.md)
  - References: [ADR 0011](/docs/adr/0011-configuration-and-secrets-management.md)

## Context

[ADR 0010](/docs/adr/0010-unified-identity-and-jwt.md) routes SSH-resolved identities through an internal, sshd-only token-exchange endpoint. If the endpoint's only protection was network binding: a misconfiguration that exposed it would let anyone mint a token for any user. Network reachability is not an authentication mechanism, and the point sharpens once sshd and the API are separate containers on a shared Docker network, where "localhost" no longer applies.

v1 closes the immediate hole with a service credential (below). This ADR also records a stronger design to adopt if and when it is worth the extra key management, so the decision is captured rather than re-derived later. It is **Proposed**: the service-credential half is adopted in v1; the asymmetric-signing half is not yet.

## Decision

Two changes, considered together:

1. **The internal exchange endpoint requires a dedicated service credential** - a high-entropy secret that sshd presents and the API verifies with a constant-time comparison, injected via the existing `_FILE` / sops mechanism ([ADR 0011](/docs/adr/0011-configuration-and-secrets-management.md)). Network binding stays only as defense-in-depth. **This part is adopted in v1**; it is the load-bearing fix and replaces a network assumption with a real credential.

2. **Switch JWTs from symmetric (HS256) to asymmetric (EdDSA).** The auth component holds the private key and is the sole signer; everything that only _verifies_ a token holds just the public key. Minting capability is then cryptographically scoped to one component instead of shared by every process that can verify. **This part is proposed, not yet adopted** - v1 keeps HS256 for fewer moving parts.

## Consequences

**Positive**

- The exchange endpoint is protected by a credential, not by where the caller sits on the network - the property originally wanted.
- Under asymmetric signing, a compromised verifier (a future read-only service, or a leaked verification key) cannot mint tokens.
- Records the upgrade path, so it is a deliberate, ready decision rather than a later scramble.

**Negative**

- Asymmetric signing adds keypair generation, distribution, and rotation. For a single-binary monolith that both signs and verifies in one process, the benefit is largely latent until there is a second, verify-only consumer - which is why it is deferred.

## Alternatives considered

- **Rely on network binding alone**. Rejected: not an authentication mechanism; the failure mode is total (mint any token) and silent.
- **mTLS between sshd and the API.** Strong, and a good thing to demonstrate, but operationally heavier than a service credential for one homelab host. Held as a further option.
