# Current architecture

## Project purpose

This service is the starting point for a future WhatsApp-based lost-dog matching product. Today it exposes only an HTTP health check so the project can grow from a tested, documented foundation.

## Current components and request flow

```text
HTTP client
    |
    | GET /health
    v
Go HTTP API (`cmd/api`)
    |
    v
Health handler (`internal/httpapi`)
```

The API process listens on port `8080`. `cmd/api` constructs the HTTP handler and starts the standard-library server. The handler routes `GET /health` and returns a JSON status response.

## Package responsibilities and dependencies

- `cmd/api` owns executable startup and process wiring.
- `internal/httpapi` owns HTTP routing and transport behavior.

Dependencies point inward from the executable to the HTTP package: `cmd/api` depends on `internal/httpapi`; both use the Go standard library. There is currently no domain or business-logic package because no domain behavior exists yet.

## Persistence and external integrations

The service is stateless. It has no persistence and makes no calls to external systems.

The current architecture explicitly has no:

- WhatsApp integration;
- database;
- image storage;
- image recognition or matching;
- authentication;
- dog catalog;
- dog report or search APIs.

## Future direction — not implemented

Future tasks may introduce report, catalog, search, and matching behavior and connect the service to WhatsApp and appropriate storage. Those decisions should be made incrementally when their requirements are known; they are not part of the current architecture.
