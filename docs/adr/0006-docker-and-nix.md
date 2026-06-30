# ADR 0006: Docker as front door, Nix for reproducibility

- **Status:** Accepted
- **Date:** 2026-06-15
- **Related Artefacts:**
  - Supplemented by: [`ADR 0011`](/docs/adr/0011-configuration-and-secrets-management.md)

## Context

The project needs to be (1) easy for any reviewer to evaluate locally and (2) deployable to my existing [NixOS homelab](https://github.com/corygyarmathy/dotfiles). Docker and Nix are both viable on their own, with different strengths.

I have invested significant time to learning and setting up Nix on my computer and homelab servers, such that the barrier to setup a develop or build process with Nix is relatively low for me.

The configuration and secrets management aspects of this are addressed in: [`ADR 0011: Configuration and secrets management`](/docs/adr/0011-configuration-and-secrets-management.md).

## Decision

Both, in parallel, with Docker as the documented quickstart:

- **Docker** (`docker compose up`) is the path advertised in the README. Universal, zero-friction, very commonly already installed. Its load-bearing job is standing up Postgres and the one-command quickstart, not packaging the binary for portability the binary does not need.
- **Nix** (`nix develop`, eventually `nix build`) provides the reproducible dev shell, deterministic builds, and the NixOS module used to deploy to the homelab. This is the production substrate.

Neither may become load-bearing for the application itself: the binary depends only on the environment-variable contract defined in [`internal/platform/config`](/internal/platform/config), so `go run ./cmd/server` against any reachable Postgres must work with no Docker and no Nix. Docker and Nix are packaging around that contract, not preconditions of it.

The final Docker image and the binary produced by `nix build` should be functionally equivalent - same compiled Go binary, just packaged differently. A stretch goal is to have Nix build the Docker image itself via `dockerTools.buildImage`, giving truly reproducible image hashes.

The usual argument for Docker, providing a consistent runtime, does not apply as the server compiles to a static binary with no runtime dependencies. The binary already runs anywhere. Compose earns its place on the **Postgres dependency**: `docker compose up` hands a reviewer a correctly-versioned, healthchecked, throwaway database with nothing installed locally.

## Consequences

**Positive**

- Reviewers can evaluate the project without learning & installing Nix, and without installing or configuring Postgres.
- Dev environment is locked across machines and CI via the flake.
- Homelab deployment is declarative and version-controlled alongside the code.

**Negative**

- Two build paths to maintain. Mitigated by them sharing the same Go source and the same `cmd/server/main.go` entrypoint.
- The Nix flake adds a directory of files that won't be relevant to reviewers who don't use Nix. Mitigated by clear README signposting.
- The clone-and-run promise depends on the schema existing after `docker compose up`. Because production runs the same binary, migrations are applied in-process at startup (advisory-locked) rather than via a separate Compose migration step; if that wiring regresses, the headline Docker justification regresses with it.

## Alternatives considered

- **Docker only.** Rejected: would mean abandoning my homelab's declarative deployment model and not demonstrating the infrastructure-as-code skills I've spent time developing.
- **Nix only.** Rejected: would shut out the majority of reviewers who don't have Nix installed and don't want to install it to evaluate a portfolio project.
- **Kubernetes cluster.** Rejected: unnecessary complexity that I do not have the physical hardware to host, nor would I want to maintain the cost of deploying on a cloud-provider.
