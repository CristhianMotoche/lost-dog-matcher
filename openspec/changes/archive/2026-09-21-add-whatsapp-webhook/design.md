## Context

See proposal.md — Why. The service is a stateless Go API using only the standard library (`cmd/api` + `internal/httpapi`), exposing `GET /health`. There is no domain logic yet, no persistence, and no external integrations. Docs/architecture.md documents only this. This change introduces the first external integration and the first protocol boundary.

## Goals / Non-Goals

**Goals:**
- A working end-to-end Meta Cloud API webhook: verified `GET` handshake, signed `POST` intake, and an outbound interactive reply sent through the Graph API.
- Keep the standard library only; no third-party dependencies.
- Preserve the current layering: our HTTP surface stays in `internal/httpapi`; Meta protocol concerns live in a new `internal/whatsapp` package.
- Keep the service stateless and synchronous for this slice.

**Non-Goals:**
- Conversation state or any follow-up intent handling beyond the two-option reply.
- Photo/media handling or image processing.
- Catalog, search, or matching behavior.
- Async delivery/queueing, retry policies beyond what the synchronous flow naturally provides.
- Message-ID deduplication beyond the stateless reply (see Risks).

## Decisions

### Meta Cloud API over legacy/gateway options
Use the current Meta Cloud API (Graph API `messages` endpoint) rather than the On-Premises API or third-party gateways (Twilio, etc.).
- Rationale: official current path, free-to-start, JSON webhook + REST send, matches "welcome message + two options" scope.
- Alternatives considered: legacy WhatsApp Business API (deprecated path, heavier ops); a gateway SDK (adds a dependency and a middleman with no requirement driving it).

### Standard library only
Use `crypto/hmac` + `crypto/sha256` for signature verification, `net/http` `http.Client` for outbound send, `encoding/json` for parsing/encoding envelopes.
- Rationale: AGENTS.md — no dependency without a clear reason; the surface is small.
- Alternatives considered: a Go WhatsApp SDK (wraps the same calls; adds version coupling for no requirement gain).

### New `internal/whatsapp` package for protocol concerns
This package owns: envelope parsing (inbound `messages` vs `statuses`), `X-Hub-Signature-256` verification, and a Graph API send client (access token, phone number ID, request encoding). `internal/httpapi` owns the routes (`/webhook/whatsapp`, GET + POST) and delegates to it.
- Rationale: keeps transport (our HTTP routes) separate from the external protocol while avoiding a premature domain layer. Domain packages should appear when domain behavior exists, not for an integration adapter.
- Alternatives considered: putting everything in `internal/httpapi` (muddles protocol with our HTTP surface); inventing a domain package now (premature — see AGENTS.md).

### Synchronous processing
Verify signature → parse envelope → on user text message, send the interactive reply → return `200` to Meta. No worker/queue.
- Rationale: one small side effect, cheap to do inline; a queue is premature for this slice.
- Alternatives considered: acknowledge `200` immediately and reply async (adds ordering/retry complexity with no volume to justify it).

### Interactive reply buttons
The reply sends a `text` body with a welcome greeting plus two interactive `reply` buttons ("Report a lost dog", "Look for one").
- Rationale: native, tappable UX on WhatsApp; the user chose this over plain-text numbered options.
- Note: interactive reply buttons are only available within the 24-hour customer-service window (or via pre-approved templates), so the button UI is best-effort within the window; this is acceptable for the MVP (see Risks).

### Fail closed
Invalid/missing signature or bad verify token → no processing, non-2xx (or minimal) response. Missing env config → process refuses to start.
- Rationale: the product will handle person-identifying chatter; defaulting to reject is the safe posture.

## Risks / Trade-offs

- [Interactive reply buttons unavailable outside the 24-hour window / without template approval] → Mitigation: reply immediately for this MVP (replies happen within the window by construction); button fallback to plain-text copy can be a follow-up.
- [Meta retries webhooks that don't return 2xx quickly; our synchronous outbound call could occasionally time out → duplicate deliveries] → Mitigation: the reply is idempotent in effect (same stateless message), so duplicates are benign today; note message-ID dedup as future work.
- [Replaying/duplicate or out-of-order events if Meta re-delivers] → Mitigation: accepted as harmless for a stateless welcome reply; revisit when the flow gains state.
- [Webhook reachable by anyone → forged/abusive payloads] → Mitigation: signature verification rejects unauthenticated payloads; verify-token handshake controls subscription.

## Migration Plan

- Deploy: this is additive — the new routes coexist with `/health`; no migration needed. Meta app webhook URL pointed at `/webhook/whatsapp`.
- Rollback: revert the route wiring/config; `/health` behavior is unchanged.

## Open Questions

- Exact welcome copy and button labels/order. Deferred: wording choice does not change the approach, specs, or tasks.