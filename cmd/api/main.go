// Command api starts the lost-dog matching HTTP service.
package main

import (
	"log"
	"net/http"

	"lostdog/internal/httpapi"
)

const address = ":8080"

func main() {
	log.Printf("API listening on %s", address)

	if err := http.ListenAndServe(address, httpapi.NewHandler()); err != nil {
		log.Fatal(err)
	}
}
