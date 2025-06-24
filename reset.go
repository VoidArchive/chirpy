package main

import (
	"context"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Forbidden", nil)
		return
	}

	err := cfg.db.DeleteAllUsers(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete users", err)
		return
	}

	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Database reset"))
}
