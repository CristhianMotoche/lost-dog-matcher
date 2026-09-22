## Purpose

Receives WhatsApp messages over the Meta Cloud API, verifies their authenticity, and starts every user conversation with a welcome message offering the two product intents: reporting a lost dog or looking for one.

## ADDED Requirements

### Requirement: Inbound webhook verification
The system SHALL expose an HTTP GET endpoint that completes the Meta Cloud API webhook verification handshake. It SHALL respond with the requested challenge only when the `hub.verify_token` query parameter matches the configured verify token; otherwise it SHALL reject the request.

#### Scenario: Successful verification
- **WHEN** the system receives a GET request with a `hub.verify_token` matching the configured verify token and a `hub.challenge` value
- **THEN** the system responds with HTTP 200 and the challenge value as the response body

#### Scenario: Verification with mismatched token
- **WHEN** the system receives a GET request with a `hub.verify_token` that does not match the configured verify token
- **THEN** the system responds with a non-200 status and no challenge value

#### Scenario: Verification with missing token
- **WHEN** the system receives a GET request without a `hub.verify_token`
- **THEN** the system responds with a non-200 status and no challenge value

### Requirement: Webhook payload signature verification
The system SHALL verify each inbound POST request using the `X-Hub-Signature-256` header, computed as HMAC-SHA256 of the raw request body signed with the configured app secret. Requests with a missing or invalid signature SHALL be rejected and not processed.

#### Scenario: Valid signature
- **WHEN** the system receives a POST request whose `X-Hub-Signature-256` matches the HMAC-SHA256 of the raw body signed with the configured app secret
- **THEN** the system processes the payload

#### Scenario: Invalid signature
- **WHEN** the system receives a POST request whose `X-Hub-Signature-256` does not match the computed signature
- **THEN** the system rejects the request and does not process the payload

#### Scenario: Missing signature
- **WHEN** the system receives a POST request with no `X-Hub-Signature-256` header
- **THEN** the system rejects the request and does not process the payload

### Requirement: Welcome reply with intent options
For an inbound text message from a user, the system SHALL send a WhatsApp message back to the sender containing a welcome greeting and interactive reply buttons for the two product intents: "Report a lost dog" and "Look for one". The welcome reply SHALL be sent only for text messages from individual users.

#### Scenario: Text message from a user
- **WHEN** the system receives a webhook payload containing a text message from a user
- **THEN** the system sends the sender a welcome message with two interactive reply buttons: "Report a lost dog" and "Look for one"

### Requirement: Ignore non-message events
The system SHALL acknowledge inbound webhook payloads that do not contain a user text message (for example, message status updates, delivery receipts, reactions, or messages containing only media) without sending any reply.

#### Scenario: Status update notification
- **WHEN** the system receives a webhook payload describing a message status update with no user text message
- **THEN** the system acknowledges the payload and sends no reply

#### Scenario: Media-only message
- **WHEN** the system receives a webhook payload containing a user message with media and no text
- **THEN** the system acknowledges the payload and sends no reply

### Requirement: Fail closed on missing configuration
The system SHALL refuse to start when any required WhatsApp configuration value (verify token, app secret, access token, phone number ID) is missing.

#### Scenario: Missing required configuration
- **WHEN** the API starts without all required WhatsApp configuration values set
- **THEN** the process exits with an error and serves no traffic