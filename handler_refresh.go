package main

import (
	"net/http"
	"time"

	"github.com/znataniel/chirpy/internal/auth"
)

func (cfg *apiConfig) refresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondJsonError(w, http.StatusUnauthorized, err, "authorization header not found")
		return
	}

	tokenRow, err := cfg.dbq.GetTokenByToken(r.Context(), refreshToken)
	if err != nil {
		respondJsonError(w, http.StatusUnauthorized, err, "token not found")
		return
	}
	if tokenRow.ExpiresAt.Before(time.Now().UTC()) {
		respondJsonError(w, http.StatusUnauthorized, err, "expired token")
		return
	}
	if tokenRow.RevokedAt.Valid {
		respondJsonError(w, http.StatusUnauthorized, err, "token has been revoked")
		return
	}

	type tokenJson struct {
		Token string `json:"token"`
	}

	newAccessToken, err := auth.MakeJWT(tokenRow.UserID, cfg.secret)
	if err != nil {
		respondJsonError(w, http.StatusInternalServerError, err, "could not create new access token")
	}

	respondJson(w, http.StatusOK, tokenJson{Token: newAccessToken})
}
