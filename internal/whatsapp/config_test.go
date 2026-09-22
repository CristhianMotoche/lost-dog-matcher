package whatsapp

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestConfigValidateValid(t *testing.T) {
	cfg := Config{
		VerifyToken:   "verify",
		AppSecret:     "secret",
		AccessToken:   "token",
		PhoneNumberID: "pid",
	}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestConfigValidateMissingValues(t *testing.T) {
	cfg := Config{}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected an error for empty config")
	}

	for _, envVar := range []string{"WHATSAPP_VERIFY_TOKEN", "WHATSAPP_APP_SECRET", "WHATSAPP_ACCESS_TOKEN", "WHATSAPP_PHONE_NUMBER_ID"} {
		if !strings.Contains(err.Error(), envVar) {
			t.Errorf("error %q does not mention %q", err, envVar)
		}
	}
}

func TestSendWelcomeComposesIntentButtons(t *testing.T) {
	var gotBody []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := Config{
		AccessToken:   "token",
		PhoneNumberID: "pid",
	}
	webhook := &Webhook{
		cfg:    cfg,
		client: NewClientWithBaseURL(cfg, server.URL),
	}

	if err := webhook.SendWelcome("15551234567"); err != nil {
		t.Fatalf("SendWelcome returned error: %v", err)
	}

	var request sendRequest
	if err := json.Unmarshal(gotBody, &request); err != nil {
		t.Fatalf("decode sent body: %v", err)
	}

	if len(request.Interactive.Action.Buttons) != 2 {
		t.Fatalf("button count = %d; want 2", len(request.Interactive.Action.Buttons))
	}
	if got := request.Interactive.Action.Buttons[0]; got.Reply.Title != "Report a lost dog" {
		t.Errorf("first button title = %q; want %q", got.Reply.Title, "Report a lost dog")
	}
	if got := request.Interactive.Action.Buttons[1]; got.Reply.Title != "Look for one" {
		t.Errorf("second button title = %q; want %q", got.Reply.Title, "Look for one")
	}
}
