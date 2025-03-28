package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/alnah/go-httpserver/internal/auth"
	"github.com/alnah/go-httpserver/internal/database"
	"github.com/google/uuid"
)

// User represents an application user.
// Fields are JSON-tagged for external representation, with internal fields (like HashedPassword) omitted.
type User struct {
	// ID is the unique identifier for the user.
	ID uuid.UUID `json:"id"`
	// Email is the user's email address.
	Email string `json:"email"`
	// IsChirpyRed indicates whether the user has an upgraded membership.
	IsChirpyRed bool `json:"is_chirpy_red"`
	// CreatedAt records when the user was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt records the last time the user was updated.
	UpdatedAt time.Time `json:"updated_at"`
	// HashedPassword stores the user's password hash (internal use only).
	HashedPassword string `json:"-"` // For internal use only.
}

// handlerUsersCreate creates a new user.
// It decodes JSON parameters, hashes the password, creates the user in the database,
// and returns the created user as JSON.
func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:          user.ID,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		},
	})
}
