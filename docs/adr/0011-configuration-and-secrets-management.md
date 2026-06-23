# ADR 0011: Configuration and secrets management

- **Status:** Accepted
- **Date:** 2026-06-17
- **Related Artefacts:**
  - References: [ADR 007](/docs/adr/0007-ssh-as-a-demo-surface.md)
  - References: [ADR 006](/docs/adr/0006-docker-and-nix.md)

## Context

The server needs runtime configuration (listen address, log level, database connection) and a small number of secrets (the JWT signing key, the database password embedded in `DATABASE_URL`). Several forces constrain how this is done:

- **The source is public.** No real secret may ever be committed, including in Nix expressions, where a naive string option lands in the world-readable `/nix/store`.
- **There is a live demo on a homelab host.** Production is real, not hypothetical, so the dev/prod boundary has to be enforced rather than trusted. (See [ADR 0007](/docs/adr/0007-ssh-as-a-demo-surface.md)).
- **Builds should be declarative and CI-friendly.** A clone-and-run quickstart is a stated goal of [ADR 0006](/docs/adr/0006-docker-and-nix.md) and CI must be able to stand the stack up with no manual secret provisioning.
- **Two deployment substrates.** The same binary is run locally via Docker Compose and in production via a NixOS module. Neither may become load-bearing for the application itself - Nix in particular must stay optional (see the [Dockerfile](/deploy/docker/Dockerfile), which builds without it). Refer to [ADR 0006](/docs/adr/0006-docker-and-nix.md).

The recurring tension is between "as declarative as possible" and "secrets". A secret that is in plaintext in version control is, by definition, leaked.

## Decision

Configuration follows [12-factor principles](https://en.wikipedia.org/wiki/Twelve-Factor_App_methodology): **the environment is the configuration interface, and the application depends only on that contract** - not on Docker, not on Nix.

Three concerns are separated deliberately:

1. **The contract** lives in [`internal/platform/config`](/internal/platform/config). It defines which variables exist, their defaults, and their validation rules. This is the only canonical source.
2. **Non-secret values** (`HTTP_ADDR`, `LOG_LEVEL`, `APP_ENV`) have safe defaults in code, so zero-config local runs work. They may be set freely in any committed file.
3. **Secret values** (`JWT_SECRET`, `DATABASE_URL`) are injected at runtime from a per-environment source and are never committed for production.

Four mechanisms implement this:

- **File-based secret injection (`_FILE` convention).** Any secret variable `KEY` may instead be supplied as `KEY_FILE`, pointing at a file whose contents are read at startup. This is the same convention used by the official Postgres image, and it is the bridge between substrates: Docker secrets mount a file, and sops-nix + systemd `LoadCredential` expose one under `$CREDENTIALS_DIRECTORY`. In both cases the secret value never enters the process environment, so it cannot leak via `/proc/<pid>/environ`.

- **A production startup guard.** When `APP_ENV=production`, the loader refuses to boot if `JWT_SECRET` is empty, shorter than 32 bytes, or still equal to the development placeholder (`dev-only-change-me`). Dev secrets crossing into production becomes a crash on boot rather than a silent vulnerability.

- **Eager, accumulating validation.** `config.Load()` is called once at startup, reports _all_ configuration problems at once, and fails the process before any dependency is wired. Misconfiguration cannot manifest as a runtime surprise.

- **Self-redacting secret type.** Secrets are held in a `Secret` type that implements `slog.LogValuer` and `fmt.Stringer` to render as `REDACTED`, with an explicit `Reveal()` accessor for the one place the value is genuinely needed. A secret cannot be logged by accident, and every real use is greppable.

Development and CI use weak, clearly-labelled placeholder secrets committed in [`deploy/docker/compose.yaml`](/deploy/docker/compose.yaml) and the [flake](flake.nix) dev shell. This is acceptable because those values protect nothing and the production guard rejects them.

## Consequences

**Positive**

- The application is portable across substrates because it depends only on the env-var contract; Docker and Nix remain interchangeable deployment details.
- No production secret ever exists in version control, the Nix store, or the process environment.
- The dev/prod boundary is enforced mechanically (the startup guard), not by discipline.
- The clone-and-run quickstart and CI both work with committed dev placeholders and require no secret provisioning.

**Negative / accepted trade-offs**

- Environment variables are flat and stringly-typed; deeply nested or structured configuration would be awkward. Acceptable - the config surface is small and expected to stay so.
- Required variables must be documented or discovered by a (deliberately loud) startup failure. Mitigated by accumulating all errors in one message and listing variables in [`compose.yaml`](/deploy/docker/compose.yaml) and the [README](/README.md). Secret rotation requires a process restart. Acceptable at this scale.
- The `_FILE` indirection is marginally more machinery than reading a plain env var. Justified by giving Docker and NixOS a single shared code path.

**Neutral**

- Once the loader supports `_FILE`, wiring the production NixOS deployment
  (phase 7) is purely a Nix exercise - `LoadCredential` plus the matching
  `*_FILE` environment variables - with no further Go changes required.

  ## Alternatives considered

| Option                                                  | Why not                                                                                                                                                                           |
| ------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Config file (TOML/YAML) committed to the repo**       | Re-introduces the secret-in-git problem, or forces a parallel secret mechanism anyway. Env vars are already the substrate-neutral interface.                                      |
| **A config library (Viper, envconfig)**                 | Overkill for ~5 fields. Hand-rolled loading is more idiomatic for this size, adds no dependency, and demonstrates the mechanics rather than hiding them.                          |
| **HashiCorp Vault / dedicated secrets manager**         | Operationally heavy for a single homelab host with one demo deployment. sops-nix already covers the production secret store. Revisit only if secret count or rotation needs grow. |
| **Baking secrets into the image or Nix store**          | The store is world-readable and the source is public. This is the specific failure mode the whole decision exists to prevent.                                                     |
| **`databaseUrl` as a plain NixOS module string option** | The embedded password would land in `/nix/store`. The module uses `databaseUrlFile` / `jwtSecretFile` (paths) instead, consistent with the `_FILE` convention above.              |
