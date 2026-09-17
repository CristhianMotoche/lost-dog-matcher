# Use devenv for local development

Status: Accepted
Date: 2026-09-13

## Context

Contributors and coding agents need a consistent Go toolchain and a shared set of commands for routine development and verification.

## Decision

Use devenv as the canonical local development environment. Enable Go through devenv's Go language support, define shared tasks for formatting, tests, vet, checks, and running the API, and configure the API as a `devenv up` process. Use devenv's git-hooks integration to run `gofmt` and `go vet` before commits.

The Go package is not explicitly versioned in `devenv.nix`; the locked nixpkgs input supplies it. This avoids pinning an older Go release while keeping the resolved environment reproducible through `devenv.lock`.

## Consequences

- Contributors use the same repository-defined development commands and resolved Go toolchain.
- Entering the devenv shell installs pre-commit checks for Go formatting and vet analysis.
- Nix and devenv are prerequisites for the canonical workflow.
- The devenv lock file must be updated deliberately when development-environment inputs change.
