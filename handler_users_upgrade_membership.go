package main

import (
	"encoding/json"
	"net/http"

	"github.com/alnah/go-httpserver/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUserUpgradeMembership(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find API key", err)
		return
	}
	if apiKey != cfg.polkaAPIKey {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate Polka API key", err)
		return
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = cfg.db.UpgradeUserMembership(r.Context(), params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't upgrade user membership", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
