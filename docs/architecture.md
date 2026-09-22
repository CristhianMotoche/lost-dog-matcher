# Current architecture

## Project purpose

This service is the starting point for a future WhatsApp-based lost-dog matching product. Today it exposes an HTTP health check and a WhatsApp webhook through which users start a conversation and choose between reporting a lost dog and looking for one.

## Current components and request flow

```text
HTTP client / Meta Cloud API
    |
    | GET /health | GET|POST /webhook/whatsapp
    v
Go HTTP API (`cmd/api`)
    |
    v
HTTP routes (`internal/httpapi`)
    |                                   |
    |                                   | delegates to
    v                                   v
Health handler                    WhatsApp protocol (`internal/whatsapp`)
                                         |
                                         | POST /v21.0/<phone-number-id>/messages
                                         v
                                    Meta Cloud API (reply)
```

The API process listens on port `8080`. `cmd/api` loads the WhatsApp configuration from the environment, validates it, and constructs the HTTP handler.

The handler routes:
- `GET /health` — returns a JSON status response.
- `GET /webhook/whatsapp` — completes the Meta Cloud API webhook verification handshake when the `hub.verify_token` matches; otherwise rejects.
- `POST /webhook/whatsapp` — verifies the `X-Hub-Signature-256` header over the raw body, and for an inbound user text message sends a welcome message back through the Meta Graph API with two interactive reply buttons ("Report a lost dog", "Look for one"). Status updates and media-only messages are acknowledged but not replied to.

## Package responsibilities and dependencies

- `cmd/api` owns executable startup, environment configuration, and process wiring.
- `internal/httpapi` owns HTTP routing, the webhook handlers, and transport behavior.
- `internal/whatsapp` owns the WhatsApp Cloud API protocol: configuration, envelope parsing, signature verification, and the outbound send client.

Dependencies point inward from the executable to the HTTP package and onward to the protocol package: `cmd/api` depends on `internal/httpapi`; `internal/httpapi` depends on `internal/whatsapp`. All use the Go standard library. There is currently no domain or business-logic package because no domain behavior exists yet.

## Persistence and external integrations

The service is stateless. It has no persistence. Its only external integration is the Meta Cloud API, used to receive webhook notifications and to send WhatsApp messages.

The current architecture explicitly has no:

- database;
- image storage;
- image recognition or matching;
- dog catalog;
- dog report or search APIs;
- conversation state or intent handling beyond the two-option welcome reply.

## Future direction — not implemented

Future tasks may introduce report, catalog, search, and matching behavior, connect the service to WhatsApp more deeply, and add appropriate storage. Those decisions should be made incrementally when their requirements are known; they are not part of the current architecture.
