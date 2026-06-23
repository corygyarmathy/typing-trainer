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

- **Docker** (`docker compose up`) is the path advertised in the README. Universal, zero-friction, very commonly already installed.
- **Nix** (`nix develop`, eventually `nix build`) provides the reproducible dev shell, deterministic builds, and the NixOS module used to deploy to the homelab.

The final Docker image and the binary produced by `nix build` should be functionally equivalent - same compiled Go binary, just packaged differently. A stretch goal is to have Nix build the Docker image itself via `dockerTools.buildImage`, giving truly reproducible image hashes.

## Consequences

**Positive**

- Reviewers can evaluate the project without learning & installing Nix.
- Dev environment is locked across machines and CI via the flake.
- Homelab deployment is declarative and version-controlled alongside the code.

**Negative**

- Two build paths to maintain. Mitigated by them sharing the same Go source and the same `cmd/server/main.go` entrypoint.
- The Nix flake adds a directory of files that won't be relevant to reviewers who don't use Nix. Mitigated by clear README signposting.

## Alternatives considered

- **Docker only.** Rejected: would mean abandoning my homelab's declarative deployment model and not demonstrating the infrastructure-as-code skills I've spent time developing.
- **Nix only.** Rejected: would shut out the majority of reviewers who don't have Nix installed and don't want to install it to evaluate a portfolio project.
- **Kubernetes cluster.** Rejected: unnecessary complexity that I do not have the physical hardware to host, nor would I want to maintain the cost of deploying on a cloud-provider.
