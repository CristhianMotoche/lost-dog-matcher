package whatsapp

import "encoding/json"

// InboundMessage is a user message reported by a WhatsApp webhook payload.
type InboundMessage struct {
	From string `json:"from"`
	ID   string `json:"id"`
	Type string `json:"type"`
	Text struct {
		Body string `json:"body"`
	} `json:"text"`
}

type inboundPayload struct {
	Entry []struct {
		Changes []struct {
			Value struct {
				Messages []InboundMessage `json:"messages"`
			} `json:"value"`
		} `json:"changes"`
	} `json:"entry"`
}

// FirstTextMessage returns the first user text message in a webhook payload,
// reporting whether one was found. Status updates and messages that carry no
// text body are ignored.
func FirstTextMessage(body []byte) (InboundMessage, bool) {
	var payload inboundPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return InboundMessage{}, false
	}

	for _, entry := range payload.Entry {
		for _, change := range entry.Changes {
			for _, message := range change.Value.Messages {
				if message.Type == "text" {
					return message, true
				}
			}
		}
	}

	return InboundMessage{}, false
}
