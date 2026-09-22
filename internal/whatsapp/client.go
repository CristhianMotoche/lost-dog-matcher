package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const defaultBaseURL = "https://graph.facebook.com"

// ReplyButton is a single interactive reply button in a WhatsApp message.
type ReplyButton struct {
	ID    string
	Title string
}

// Client sends messages through the WhatsApp Cloud API.
type Client struct {
	baseURL     string
	accessToken string
	phoneID     string
	httpClient  *http.Client
}

// NewClient builds a Client for the given configuration.
func NewClient(cfg Config) *Client {
	return NewClientWithBaseURL(cfg, defaultBaseURL)
}

// NewClientWithBaseURL builds a Client targeting a custom base URL, primarily
// for tests.
func NewClientWithBaseURL(cfg Config, baseURL string) *Client {
	return &Client{
		baseURL:     baseURL,
		accessToken: cfg.AccessToken,
		phoneID:     cfg.PhoneNumberID,
		httpClient:  &http.Client{},
	}
}

// SendInteractive sends an interactive message with reply buttons to a user.
func (c *Client) SendInteractive(to string, header, body string, buttons []ReplyButton) error {
	request := buildSendRequest(to, header, body, buttons)
	buf, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("marshal send request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/v21.0/"+c.phoneID+"/messages", bytes.NewReader(buf))
	if err != nil {
		return fmt.Errorf("build send request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send whatsapp message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("whatsapp send failed: %s: %s", resp.Status, respBody)
	}

	return nil
}

type sendRequest struct {
	MessagingProduct string          `json:"messaging_product"`
	RecipientType    string          `json:"recipient_type"`
	To               string          `json:"to"`
	Type             string          `json:"type"`
	Interactive      interactiveBody `json:"interactive"`
}

type interactiveBody struct {
	Type   string      `json:"type"`
	Header messageText `json:"header"`
	Body   messageText `json:"body"`
	Action action      `json:"action"`
}

type messageText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type action struct {
	Buttons []button `json:"buttons"`
}

type button struct {
	Type  string      `json:"type"`
	Reply replyButton `json:"reply"`
}

type replyButton struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func buildSendRequest(to, header, body string, buttons []ReplyButton) sendRequest {
	request := sendRequest{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               to,
		Type:             "interactive",
		Interactive: interactiveBody{
			Type:   "button",
			Header: messageText{Type: "text", Text: header},
			Body:   messageText{Type: "text", Text: body},
		},
	}
	for _, b := range buttons {
		request.Interactive.Action.Buttons = append(request.Interactive.Action.Buttons, button{
			Type:  "reply",
			Reply: replyButton{ID: b.ID, Title: b.Title},
		})
	}
	return request
}
