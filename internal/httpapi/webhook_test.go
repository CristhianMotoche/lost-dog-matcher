package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type stubWhatsApp struct {
	token      string
	verifyOK   bool
	sendErr    error
	welcomedTo []string
}

func (s *stubWhatsApp) VerifyToken() string { return s.token }

func (s *stubWhatsApp) VerifySignature(body []byte, signature string) bool {
	return s.verifyOK
}

func (s *stubWhatsApp) SendWelcome(to string) error {
	if s.sendErr != nil {
		return s.sendErr
	}
	s.welcomedTo = append(s.welcomedTo, to)
	return nil
}

func webhookRequest(t *testing.T, method, target string, body string) *http.Request {
	t.Helper()
	request := httptest.NewRequest(method, target, strings.NewReader(body))
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("X-Hub-Signature-256", "sha256=test")
	}
	return request
}

func TestWebhookVerificationSuccess(t *testing.T) {
	wa := &stubWhatsApp{token: "verify-token"}
	handler := NewHandlerWithWhatsApp(wa)

	request := httptest.NewRequest(http.MethodGet, "/webhook/whatsapp?hub.verify_token=verify-token&hub.challenge=challenge-123&hub.mode=subscribe", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status code = %d; want %d", response.Code, http.StatusOK)
	}
	if body := response.Body.String(); body != "challenge-123" {
		t.Errorf("body = %q; want %q", body, "challenge-123")
	}
}

func TestWebhookVerificationTokenMismatch(t *testing.T) {
	wa := &stubWhatsApp{token: "verify-token"}
	handler := NewHandlerWithWhatsApp(wa)

	request := httptest.NewRequest(http.MethodGet, "/webhook/whatsapp?hub.verify_token=wrong&hub.challenge=challenge-123", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("status code = %d; want %d", response.Code, http.StatusForbidden)
	}
	if body := response.Body.String(); strings.Contains(body, "challenge-123") {
		t.Errorf("body %q must not contain the challenge", body)
	}
}

func TestWebhookVerificationMissingToken(t *testing.T) {
	wa := &stubWhatsApp{token: "verify-token"}
	handler := NewHandlerWithWhatsApp(wa)

	request := httptest.NewRequest(http.MethodGet, "/webhook/whatsapp?hub.challenge=challenge-123", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("status code = %d; want %d", response.Code, http.StatusForbidden)
	}
}

func TestWebhookEventRepliesToUserText(t *testing.T) {
	wa := &stubWhatsApp{verifyOK: true}
	handler := NewHandlerWithWhatsApp(wa)

	payload := `{
		"entry": [{
			"changes": [{
				"value": {
					"messages": [{"from": "15551234567", "id": "wamid.1", "type": "text", "text": {"body": "hello"}}]
				}
			}]
		}]
	}`
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, webhookRequest(t, http.MethodPost, "/webhook/whatsapp", payload))

	if response.Code != http.StatusOK {
		t.Fatalf("status code = %d; want %d", response.Code, http.StatusOK)
	}
	if len(wa.welcomedTo) != 1 || wa.welcomedTo[0] != "15551234567" {
		t.Errorf("welcomedTo = %v; want [15551234567]", wa.welcomedTo)
	}
}

func TestWebhookEventIgnoresStatusUpdate(t *testing.T) {
	wa := &stubWhatsApp{verifyOK: true}
	handler := NewHandlerWithWhatsApp(wa)

	payload := `{
		"entry": [{
			"changes": [{
				"value": {
					"statuses": [{"id": "wamid.status", "status": "delivered", "recipient_id": "15551234567"}]
				}
			}]
		}]
	}`
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, webhookRequest(t, http.MethodPost, "/webhook/whatsapp", payload))

	if response.Code != http.StatusOK {
		t.Fatalf("status code = %d; want %d", response.Code, http.StatusOK)
	}
	if len(wa.welcomedTo) != 0 {
		t.Errorf("welcomedTo = %v; want none", wa.welcomedTo)
	}
}

func TestWebhookEventIgnoresMediaOnlyMessage(t *testing.T) {
	wa := &stubWhatsApp{verifyOK: true}
	handler := NewHandlerWithWhatsApp(wa)

	payload := `{
		"entry": [{
			"changes": [{
				"value": {
					"messages": [{"from": "15551234567", "id": "wamid.1", "type": "image", "image": {"id": "media-1"}}]
				}
			}]
		}]
	}`
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, webhookRequest(t, http.MethodPost, "/webhook/whatsapp", payload))

	if response.Code != http.StatusOK {
		t.Fatalf("status code = %d; want %d", response.Code, http.StatusOK)
	}
	if len(wa.welcomedTo) != 0 {
		t.Errorf("welcomedTo = %v; want none", wa.welcomedTo)
	}
}

func TestWebhookEventRejectsInvalidSignature(t *testing.T) {
	wa := &stubWhatsApp{verifyOK: false}
	handler := NewHandlerWithWhatsApp(wa)

	payload := `{
		"entry": [{
			"changes": [{
				"value": {
					"messages": [{"from": "15551234567", "id": "wamid.1", "type": "text", "text": {"body": "hello"}}]
				}
			}]
		}]
	}`
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, webhookRequest(t, http.MethodPost, "/webhook/whatsapp", payload))

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status code = %d; want %d", response.Code, http.StatusUnauthorized)
	}
	if len(wa.welcomedTo) != 0 {
		t.Errorf("welcomedTo = %v; want none", wa.welcomedTo)
	}
}

func TestWebhookEventRejectsUnsupportedMethod(t *testing.T) {
	wa := &stubWhatsApp{verifyOK: true}
	handler := NewHandlerWithWhatsApp(wa)

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodPut, "/webhook/whatsapp", nil))

	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status code = %d; want %d", response.Code, http.StatusMethodNotAllowed)
	}
	if allow := response.Header().Get("Allow"); allow != "GET, POST" {
		t.Errorf("Allow = %q; want %q", allow, "GET, POST")
	}
}

func TestWebhookNotRegisteredWithoutWhatsApp(t *testing.T) {
	response := httptest.NewRecorder()
	NewHandler().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/webhook/whatsapp", nil))

	if response.Code != http.StatusNotFound {
		t.Fatalf("status code = %d; want %d", response.Code, http.StatusNotFound)
	}
}

func TestHealthStillServedWithWebhookWired(t *testing.T) {
	wa := &stubWhatsApp{verifyOK: true}
	handler := NewHandlerWithWhatsApp(wa)

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/health", nil))

	if response.Code != http.StatusOK {
		t.Fatalf("status code = %d; want %d", response.Code, http.StatusOK)
	}
}
