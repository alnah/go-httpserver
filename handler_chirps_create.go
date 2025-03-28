package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/alnah/go-httpserver/internal/auth"
	"github.com/alnah/go-httpserver/internal/database"
	"github.com/google/uuid"
)

// Chirp represents a short message (chirp) posted by a user.
type Chirp struct {
	// ID is the unique identifier of the chirp.
	ID uuid.UUID `json:"id"`
	// CreatedAt is the timestamp when the chirp was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the timestamp when the chirp was last updated.
	UpdatedAt time.Time `json:"updated_at"`
	// UserID is the identifier of the user who posted the chirp.
	UserID uuid.UUID `json:"user_id"`
	// Body is the content of the chirp.
	Body string `json:"body"`
}

// handlerChirpsCreate creates a new chirp.
// It validates the user's JWT, decodes the chirp content, cleans it by filtering banned words,
// and inserts the new chirp into the database.
func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

// validateChirp checks that the chirp's body does not exceed the maximum allowed length
// and filters out any banned words. It returns the cleaned chirp or an error.
func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(body, badWords)
	return cleaned, nil
}

// getCleanedBody replaces banned words in the chirp body with asterisks.
func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
