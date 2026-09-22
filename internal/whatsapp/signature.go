package whatsapp

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

const signaturePrefix = "sha256="

// VerifySignature reports whether signature matches the HMAC-SHA256 of body
// computed with the app secret, as expected in the X-Hub-Signature-256 header.
func VerifySignature(appSecret string, body []byte, signature string) bool {
	if !strings.HasPrefix(signature, signaturePrefix) {
		return false
	}

	mac := hmac.New(sha256.New, []byte(appSecret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	provided := strings.TrimPrefix(signature, signaturePrefix)

	return hmac.Equal([]byte(expected), []byte(provided))
}
