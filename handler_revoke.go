package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/znataniel/chirpy/internal/auth"
	"github.com/znataniel/chirpy/internal/database"
)

func (cfg *apiConfig) revoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondJsonError(w, http.StatusUnauthorized, err, "authorization header not found")
		return
	}

	if err := cfg.dbq.RevokeToken(r.Context(), database.RevokeTokenParams{
		Token:     refreshToken,
		RevokedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
		UpdatedAt: time.Now().UTC(),
	}); err != nil {
		respondJsonError(w, http.StatusInternalServerError, err, "could not revoke token in db")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
