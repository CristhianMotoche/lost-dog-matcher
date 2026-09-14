// Package httpapi contains the service's HTTP transport.
package httpapi

import "net/http"

// NewHandler builds the HTTP handler for the API.
func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", health)

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
