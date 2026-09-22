package whatsapp

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func signBody(appSecret string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(appSecret))
	mac.Write(body)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func TestVerifySignatureValid(t *testing.T) {
	const appSecret = "s3cret"
	body := []byte(`{"object":"whatsapp_business_account"}`)

	if !VerifySignature(appSecret, body, signBody(appSecret, body)) {
		t.Error("expected signature to verify")
	}
}

func TestVerifySignatureInvalid(t *testing.T) {
	const appSecret = "s3cret"
	body := []byte(`{"object":"whatsapp_business_account"}`)

	if VerifySignature(appSecret, body, signBody(appSecret, []byte("tampered"))) {
		t.Error("expected signature to be rejected")
	}
}

func TestVerifySignatureMissing(t *testing.T) {
	const appSecret = "s3cret"
	body := []byte(`{"object":"whatsapp_business_account"}`)

	if VerifySignature(appSecret, body, "") {
		t.Error("expected missing signature to be rejected")
	}
}

func TestVerifySignatureWrongSecret(t *testing.T) {
	body := []byte(`{"object":"whatsapp_business_account"}`)

	if VerifySignature("other-secret", body, signBody("s3cret", body)) {
		t.Error("expected signature from another secret to be rejected")
	}
}
