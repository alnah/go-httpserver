package main

import "net/http"

// handlerReadiness is a simple health-check endpoint.
// It returns a plain text "OK" message along with an HTTP 200 status.
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(http.StatusText(http.StatusOK)))
}
