// Command api starts the lost-dog matching HTTP service.
package main

import (
	"log"
	"net/http"

	"lostdog/internal/httpapi"
	"lostdog/internal/whatsapp"
)

const address = ":8080"

func main() {
	cfg := whatsapp.FromEnv()
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	webhook := whatsapp.NewWebhook(cfg)

	log.Printf("API listening on %s", address)

	if err := http.ListenAndServe(address, httpapi.NewHandlerWithWhatsApp(webhook)); err != nil {
		log.Fatal(err)
	}
}
