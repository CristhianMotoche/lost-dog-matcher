package whatsapp

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSendInteractive(t *testing.T) {
	var gotPath string
	var gotAuth string
	var gotBody []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotAuth = r.Header.Get("Authorization")
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := Config{
		AccessToken:   "access-token",
		PhoneNumberID: "PHONE_NUMBER_ID",
	}
	client := NewClientWithBaseURL(cfg, server.URL)

	err := client.SendInteractive("15551234567", "Welcome header", "Welcome body", []ReplyButton{
		{ID: "report_lost_dog", Title: "Report a lost dog"},
		{ID: "look_for_one", Title: "Look for one"},
	})
	if err != nil {
		t.Fatalf("SendInteractive returned error: %v", err)
	}

	if gotPath != "/v21.0/PHONE_NUMBER_ID/messages" {
		t.Errorf("path = %q; want %q", gotPath, "/v21.0/PHONE_NUMBER_ID/messages")
	}
	if gotAuth != "Bearer access-token" {
		t.Errorf("Authorization = %q; want %q", gotAuth, "Bearer access-token")
	}

	var request sendRequest
	if err := json.Unmarshal(gotBody, &request); err != nil {
		t.Fatalf("decode sent body: %v", err)
	}

	if request.MessagingProduct != "whatsapp" {
		t.Errorf("MessagingProduct = %q; want %q", request.MessagingProduct, "whatsapp")
	}
	if request.RecipientType != "individual" {
		t.Errorf("RecipientType = %q; want %q", request.RecipientType, "individual")
	}
	if request.To != "15551234567" {
		t.Errorf("To = %q; want %q", request.To, "15551234567")
	}
	if request.Type != "interactive" {
		t.Errorf("Type = %q; want %q", request.Type, "interactive")
	}
	if request.Interactive.Type != "button" {
		t.Errorf("Interactive.Type = %q; want %q", request.Interactive.Type, "button")
	}
	if request.Interactive.Header.Text != "Welcome header" {
		t.Errorf("Interactive.Header.Text = %q; want %q", request.Interactive.Header.Text, "Welcome header")
	}
	if request.Interactive.Body.Text != "Welcome body" {
		t.Errorf("Interactive.Body.Text = %q; want %q", request.Interactive.Body.Text, "Welcome body")
	}

	if len(request.Interactive.Action.Buttons) != 2 {
		t.Fatalf("button count = %d; want 2", len(request.Interactive.Action.Buttons))
	}
	if got := request.Interactive.Action.Buttons[0]; got.Type != "reply" || got.Reply.ID != "report_lost_dog" || got.Reply.Title != "Report a lost dog" {
		t.Errorf("first button = %+v", got)
	}
	if got := request.Interactive.Action.Buttons[1]; got.Type != "reply" || got.Reply.ID != "look_for_one" || got.Reply.Title != "Look for one" {
		t.Errorf("second button = %+v", got)
	}
}

func TestSendInteractiveReportsGraphAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":{"message":"boom"}}`, http.StatusBadRequest)
	}))
	defer server.Close()

	client := NewClientWithBaseURL(Config{PhoneNumberID: "PID"}, server.URL)

	err := client.SendInteractive("15551234567", "h", "b", []ReplyButton{{ID: "x", Title: "X"}})
	if err == nil {
		t.Fatal("expected an error for a non-2xx Graph API response")
	}
	if !strings.Contains(err.Error(), "boom") {
		t.Errorf("error = %q; want it to include the Graph API error body", err)
	}
}
