// Package whatsapp integrates with the Meta Cloud API for WhatsApp.
package whatsapp

import (
	"fmt"
	"os"
	"strings"
)

// Config holds the credentials and identifiers needed to talk to the
// WhatsApp Cloud API.
type Config struct {
	VerifyToken   string
	AppSecret     string
	AccessToken   string
	PhoneNumberID string
}

// FromEnv loads the WhatsApp configuration from the WHATSAPP_* environment
// variables.
func FromEnv() Config {
	return Config{
		VerifyToken:   os.Getenv("WHATSAPP_VERIFY_TOKEN"),
		AppSecret:     os.Getenv("WHATSAPP_APP_SECRET"),
		AccessToken:   os.Getenv("WHATSAPP_ACCESS_TOKEN"),
		PhoneNumberID: os.Getenv("WHATSAPP_PHONE_NUMBER_ID"),
	}
}

// Validate reports whether all required configuration values are present.
func (c Config) Validate() error {
	var missing []string
	if c.VerifyToken == "" {
		missing = append(missing, "WHATSAPP_VERIFY_TOKEN")
	}
	if c.AppSecret == "" {
		missing = append(missing, "WHATSAPP_APP_SECRET")
	}
	if c.AccessToken == "" {
		missing = append(missing, "WHATSAPP_ACCESS_TOKEN")
	}
	if c.PhoneNumberID == "" {
		missing = append(missing, "WHATSAPP_PHONE_NUMBER_ID")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required WhatsApp configuration: %s", strings.Join(missing, ", "))
	}
	return nil
}
