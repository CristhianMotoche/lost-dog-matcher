// Package httpapi contains the service's HTTP transport.
package httpapi

import (
	"io"
	"net/http"

	"lostdog/internal/whatsapp"
)

const webhookPath = "/webhook/whatsapp"

// WhatsApp is the subset of the WhatsApp integration the HTTP layer uses.
type WhatsApp interface {
	VerifyToken() string
	VerifySignature(body []byte, signature string) bool
	SendWelcome(to string) error
}

// NewHandler builds the HTTP handler for the API.
func NewHandler() http.Handler {
	return buildHandler(nil)
}

// NewHandlerWithWhatsApp builds the HTTP handler for the API, including the
// WhatsApp webhook routes.
func NewHandlerWithWhatsApp(wa WhatsApp) http.Handler {
	return buildHandler(wa)
}

func buildHandler(wa WhatsApp) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", health)

	if wa != nil {
		mux.HandleFunc(webhookPath, handleWebhook(wa))
	}

	return mux
}

func health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"status\":\"ok\"}\n"))
}

func handleWebhook(wa WhatsApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleWebhookVerify(w, r, wa)
		case http.MethodPost:
			handleWebhookEvent(w, r, wa)
		default:
			w.Header().Set("Allow", "GET, POST")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func handleWebhookVerify(w http.ResponseWriter, r *http.Request, wa WhatsApp) {
	query := r.URL.Query()
	if query.Get("hub.verify_token") != wa.VerifyToken() {
		http.Error(w, "verification token mismatch", http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, query.Get("hub.challenge"))
}

func handleWebhookEvent(w http.ResponseWriter, r *http.Request, wa WhatsApp) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	if !wa.VerifySignature(body, r.Header.Get("X-Hub-Signature-256")) {
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}

	if message, ok := whatsapp.FirstTextMessage(body); ok {
		if err := wa.SendWelcome(message.From); err != nil {
			http.Error(w, "failed to send reply", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
