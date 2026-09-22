## Why

The product is delivered entirely through WhatsApp, yet the service currently exposes only a `GET /health` endpoint. Users need a real entry point: an inbound WhatsApp message should greet the sender and offer the two intents that define the product — reporting a lost dog or looking for one. This change establishes that entry point so the conversation flow can grow incrementally from it.

## What Changes

- Add a WhatsApp webhook to the existing Go HTTP API for the Meta Cloud API:
  - `GET` route for the webhook verification handshake.
  - `POST` route for inbound message notifications.
- Validate inbound webhook payloads using the `X-Hub-Signature-256` HMAC signature.
- On an inbound text message from a user, reply through the Meta send-message API with interactive reply buttons offering two options: "Report a lost dog" and "Look for one".
- Add a small outgoing HTTP client for the Meta Graph API (`/v21.0/<phone-number-id>/messages`).
- Configure the webhook via environment variables: verify token, app secret, access token, and phone number ID.
- Ignore inbound events that are not the relevant text message type (e.g., status updates, delivery receipts, media).

Non-goals for this change:
- No image processing or photo handling.
- No actual report cataloging, search, or matching — selecting an option only acknowledges the choice; the follow-up flow is future work.
- No interactive message state or conversation memory.

## Capabilities

### New Capabilities
- `whatsapp/webhook`: receiving WhatsApp messages via the Meta Cloud API, verifying them, and replying with the introductory intent options.

### Modified Capabilities
- None. This is the first capability; there are no existing specs to change.

## Impact

- `cmd/api`: reads new configuration from the environment and passes it to the handler wiring.
- `internal/httpapi`: new webhook routes and request handling; keeps transport specifics in the HTTP layer.
- New `internal/whatsapp` package: Meta Cloud API message parsing, signature verification, and outbound send client. Introduced because the WhatsApp transport and protocol now become a real concern distinct from HTTP.
- Dependencies: the Go standard library only — signatures and HTTP are already covered. No new third-party packages.
- Documentation: `docs/architecture.md` updated to reflect the new component and request flow.
- Operational notes: the service remains stateless; no external storage is introduced.