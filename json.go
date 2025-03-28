package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// respondWithError logs the provided error (if any) and sends a JSON error response
// with the specified HTTP status code and message.
func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

// respondWithJSON sends a JSON response with the provided payload and HTTP status code.
// It sets the Content-Type header to "application/json" and handles marshalling errors.
func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	_, _ = w.Write(dat)
}
