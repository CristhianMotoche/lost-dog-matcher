package whatsapp

// Webhook exposes the protocol operations the HTTP layer needs to serve the
// WhatsApp Cloud API webhook.
type Webhook struct {
	cfg    Config
	client *Client
}

// NewWebhook wires a Webhook for the given configuration.
func NewWebhook(cfg Config) *Webhook {
	return &Webhook{cfg: cfg, client: NewClient(cfg)}
}

// VerifyToken returns the configured verify token.
func (w *Webhook) VerifyToken() string {
	return w.cfg.VerifyToken
}

// VerifySignature reports whether the given X-Hub-Signature-256 header value
// is valid for the request body.
func (w *Webhook) VerifySignature(body []byte, signature string) bool {
	return VerifySignature(w.cfg.AppSecret, body, signature)
}

// SendWelcome sends the introductory welcome message with the two intent
// buttons to the given WhatsApp user.
func (w *Webhook) SendWelcome(to string) error {
	const (
		header = "Welcome to Lost Dog Matcher"
		body   = "Let's get started. What can we help you with?"
	)

	buttons := []ReplyButton{
		{ID: "report_lost_dog", Title: "Report a lost dog"},
		{ID: "look_for_one", Title: "Look for one"},
	}

	return w.client.SendInteractive(to, header, body, buttons)
}
