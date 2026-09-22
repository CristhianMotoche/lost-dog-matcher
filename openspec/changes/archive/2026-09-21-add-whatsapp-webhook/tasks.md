## 1. Configuration and wiring

- [x] 1.1 Add a `Config` struct in `internal/whatsapp` holding verify token, app secret, access token, and phone number ID
- [x] 1.2 Load WhatsApp config from env vars (`WHATSAPP_VERIFY_TOKEN`, `WHATSAPP_APP_SECRET`, `WHATSAPP_ACCESS_TOKEN`, `WHATSAPP_PHONE_NUMBER_ID`) in `cmd/api`
- [x] 1.3 Fail fast at startup when any required WhatsApp config value is missing
- [x] 1.4 Wire the webhook: add `GET` and `POST` `/webhook/whatsapp` routes in `internal/httpapi`, delegating to `internal/whatsapp`

## 2. Webhook verification handshake

- [x] 2.1 Implement GET verification handler: echo `hub.challenge` when `hub.verify_token` matches, otherwise reject
- [x] 2.2 Add tests for successful, mismatched-token, and missing-token verification

## 3. Signature verification

- [x] 3.1 Implement `X-Hub-Signature-256` verification using HMAC-SHA256 over the raw request body with the app secret
- [x] 3.2 Reject requests with missing or invalid signatures before any processing
- [x] 3.3 Add tests for valid, invalid, and missing signatures

## 4. Message handling and reply

- [x] 4.1 Parse inbound payloads, distinguishing user text messages from status updates and media-only messages
- [x] 4.2 Acknowledge and ignore non-text-message events without replying
- [x] 4.3 Implement the Graph API send client (`POST /v21.0/<phone-number-id>/messages`) using the access token
- [x] 4.4 Compose and send the welcome message with the two interactive reply buttons ("Report a lost dog", "Look for one") for user text messages
- [x] 4.5 Add tests for envelope parsing, event ignoring, reply payload composition, and the send client (using `httptest` against a stub server)

## 5. Verification and documentation

- [x] 5.1 Update `docs/architecture.md` with the new component, routes, and request flow
- [x] 5.2 Run `go fmt ./...`, `go vet ./...`, and `go test ./...`; run the full `devenv tasks run app:check` and confirm green