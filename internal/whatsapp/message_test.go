package whatsapp

import "testing"

func textMessagePayload(t *testing.T, messageType, body string) []byte {
	t.Helper()
	return []byte(`{
		"object": "whatsapp_business_account",
		"entry": [{
			"id": "WHATSAPP_BUSINESS_ACCOUNT_ID",
			"changes": [{
				"value": {
					"messaging_product": "whatsapp",
					"metadata": {
						"display_phone_number": "15550000000",
						"phone_number_id": "PHONE_NUMBER_ID"
					},
					"contacts": [{
						"profile": {"name": "Test User"},
						"wa_id": "15551234567"
					}],
					"messages": [{
						"from": "15551234567",
						"id": "wamid.test",
						"timestamp": "1694907770",
						"type": "` + messageType + `",
						"` + messageType + `": {"body": "` + body + `"}
					}]
				},
				"field": "messages"
			}]
		}]
	}`)
}

func statusUpdatePayload() []byte {
	return []byte(`{
		"object": "whatsapp_business_account",
		"entry": [{
			"id": "WHATSAPP_BUSINESS_ACCOUNT_ID",
			"changes": [{
				"value": {
					"messaging_product": "whatsapp",
					"metadata": {
						"display_phone_number": "15550000000",
						"phone_number_id": "PHONE_NUMBER_ID"
					},
					"statuses": [{
						"id": "wamid.status",
						"status": "delivered",
						"timestamp": "1694907770",
						"recipient_id": "15551234567"
					}]
				},
				"field": "messages"
			}]
		}]
	}`)
}

func TestFirstTextMessageFindsUserText(t *testing.T) {
	message, ok := FirstTextMessage(textMessagePayload(t, "text", "hello"))

	if !ok {
		t.Fatal("expected a text message to be found")
	}
	if message.From != "15551234567" {
		t.Errorf("From = %q; want %q", message.From, "15551234567")
	}
	if message.Type != "text" {
		t.Errorf("Type = %q; want %q", message.Type, "text")
	}
	if message.Text.Body != "hello" {
		t.Errorf("Text.Body = %q; want %q", message.Text.Body, "hello")
	}
}

func TestFirstTextMessageSkipsNonTextMessage(t *testing.T) {
	if _, ok := FirstTextMessage(textMessagePayload(t, "image", "")); ok {
		t.Error("expected no text message for a media-only payload")
	}
}

func TestFirstTextMessageSkipsStatusUpdate(t *testing.T) {
	if _, ok := FirstTextMessage(statusUpdatePayload()); ok {
		t.Error("expected no text message for a status update payload")
	}
}

func TestFirstTextMessageHandlesMalformedPayload(t *testing.T) {
	if _, ok := FirstTextMessage([]byte("not json")); ok {
		t.Error("expected no text message for malformed payload")
	}
}
