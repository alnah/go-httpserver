package main

import "net/http"

// handlerReset resets the file server hit counter and resets the database to its initial state.
// This operation is only allowed when running in a development environment.
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("Reset is only allowed in dev environment."))
		return
	}

	cfg.fileserverHits.Store(0)
	_ = cfg.db.Reset(r.Context())
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Hits reset to 0 and database reset to initial state."))
}
