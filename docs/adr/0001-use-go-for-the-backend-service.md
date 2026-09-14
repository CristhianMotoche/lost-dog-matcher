# Use Go for the backend service

Status: Accepted
Date: 2026-09-13

## Context

The project needs a readable, maintainable backend foundation for an HTTP API that may later support a WhatsApp-based lost-dog matching service. The initial service has modest requirements and should be easy to learn and operate.

## Decision

Use Go for the backend service. Begin with the standard library and a conventional `cmd` and `internal` package layout.

## Consequences

- The service has a small runtime footprint and straightforward build and test commands.
- Contributors can use Go's standard formatting, testing, and vet tooling.
- The implementation should favor clear, idiomatic Go so it remains approachable to developers learning the language.
- The team must maintain Go knowledge and tooling.
